package webapp

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

func (r *WebAppRouter) budgetPlanningHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "budget_planning")

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}
	data["UserID"] = userID

	// Handle POST request (save budget or copy)
	if req.Method == http.MethodPost {
		if req.FormValue("action") == "copy" {
			r.handleBudgetCopyPOST(w, req, userID)
		} else {
			r.handleBudgetPlanningPOST(w, req, userID)
		}
		return
	}

	// Parse month parameter (default to next month)
	monthStart := getNextMonth()
	if monthParam := req.URL.Query().Get("month"); monthParam != "" {
		if timestamp, err := strconv.ParseInt(monthParam, 10, 64); err == nil {
			monthStart = time.Unix(timestamp, 0)
		}
	}
	// Ensure we use first day of month
	monthStart = time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())
	data["MonthStart"] = monthStart

	// Get all accounts and filter expense accounts
	accounts, err := r.db.GetAccounts(userID)
	if err != nil {
		r.logger.Error("Failed to get accounts", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	expenseAccounts := make([]goserver.Account, 0)
	for _, account := range accounts {
		if account.Type == "expense" {
			expenseAccounts = append(expenseAccounts, account)
		}
	}
	data["ExpenseAccounts"] = expenseAccounts

	// Get existing budget items for this month
	budgetItems, err := r.budgetService.ListMonthlyBudget(req.Context(), userID, monthStart)
	if err != nil {
		r.logger.Error("Failed to get budget items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["BudgetItems"] = budgetItems

	if err := tmpl.ExecuteTemplate(w, "budget_planning.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getNextMonth returns the first day of next month
func getNextMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
}

func (r *WebAppRouter) handleBudgetPlanningPOST(w http.ResponseWriter, req *http.Request, userID string) {
	// Parse form data
	if err := req.ParseForm(); err != nil {
		r.logger.Error("Failed to parse form", "error", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse month
	monthStart := getNextMonth()
	if monthParam := req.FormValue("month"); monthParam != "" {
		if timestamp, err := strconv.ParseInt(monthParam, 10, 64); err == nil {
			monthStart = time.Unix(timestamp, 0)
		}
	}
	monthStart = time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())

	// Build budget entries from form data
	var budgetEntries []goserver.BudgetItemNoId
	for key, values := range req.PostForm {
		if len(values) > 0 && values[0] != "" && key != "month" {
			// key should be account ID, value should be amount
			if amount, err := strconv.ParseFloat(values[0], 64); err == nil && amount > 0 {
				budgetEntries = append(budgetEntries, goserver.BudgetItemNoId{
					AccountId:   key,
					Amount:      amount,
					Date:        monthStart,
					Description: "Monthly budget",
				})
			}
		}
	}

	// Save budget
	if err := r.budgetService.SaveMonthlyBudget(req.Context(), userID, monthStart, budgetEntries); err != nil {
		r.logger.Error("Failed to save budget", "error", err)
		// Redirect back with error
		redirectURL := req.URL.Path + "?month=" + strconv.FormatInt(monthStart.Unix(), 10) + "&error=" + err.Error()
		http.Redirect(w, req, redirectURL, http.StatusSeeOther)
		return
	}

	// Redirect back with success
	redirectURL := req.URL.Path + "?month=" + strconv.FormatInt(monthStart.Unix(), 10) +
		"&success=Budget saved successfully"
	http.Redirect(w, req, redirectURL, http.StatusSeeOther)
}

func (r *WebAppRouter) handleBudgetCopyPOST(w http.ResponseWriter, req *http.Request, userID string) {
	// Parse form data
	if err := req.ParseForm(); err != nil {
		r.logger.Error("Failed to parse form", "error", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse target month
	toMonthStart := getNextMonth()
	if monthParam := req.FormValue("month"); monthParam != "" {
		if timestamp, err := strconv.ParseInt(monthParam, 10, 64); err == nil {
			toMonthStart = time.Unix(timestamp, 0)
		}
	}
	toMonthStart = time.Date(toMonthStart.Year(), toMonthStart.Month(), 1, 0, 0, 0, 0, toMonthStart.Location())

	// Default from month is previous month
	fromMonthStart := toMonthStart.AddDate(0, -1, 0)

	// Copy budget
	count, err := r.budgetService.CopyFromPreviousMonth(req.Context(), userID, fromMonthStart, toMonthStart)
	if err != nil {
		r.logger.Error("Failed to copy budget", "error", err)
		// Redirect back with error
		redirectURL := req.URL.Path + "?month=" + strconv.FormatInt(toMonthStart.Unix(), 10) + "&error=" + err.Error()
		http.Redirect(w, req, redirectURL, http.StatusSeeOther)
		return
	}

	// Redirect back with success
	successMsg := fmt.Sprintf("Copied %d budget items from %s", count, fromMonthStart.Format("January 2006"))
	redirectURL := req.URL.Path + "?month=" + strconv.FormatInt(toMonthStart.Unix(), 10) + "&success=" + successMsg
	http.Redirect(w, req, redirectURL, http.StatusSeeOther)
}
