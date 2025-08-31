package webapp

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

// normalizeMonth returns the first day of the given month at 00:00 in the same location
func normalizeMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, t.Location())
}

// budgetUnifiedHandler serves GET and POST for /web/budget
// GET: renders unified planning+comparison page
// POST: saves plan or copies from previous month depending on action
func (r *WebAppRouter) budgetUnifiedHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		r.logger.Error("Failed to load templates", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := utils.CreateTemplateData(req, "budget")

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		// ValidateUserID should have already handled the response (login page or error)
		return
	}
	data["UserID"] = userID

	if req.Method == http.MethodPost {
		r.handleUnifiedBudgetPOST(w, req, userID)
		return
	}

	// GET flow
	monthStart := normalizeMonth(time.Now()) // default to current month
	if monthParam := req.URL.Query().Get("month"); monthParam != "" {
		if ts, err := strconv.ParseInt(monthParam, 10, 64); err == nil {
			monthStart = time.Unix(ts, 0)
		}
	}
	monthStart = normalizeMonth(monthStart)
	data["MonthStart"] = monthStart
	data["PrevMonth"] = monthStart.AddDate(0, -1, 0)
	data["NextMonth"] = monthStart.AddDate(0, 1, 0)

	// Get all accounts and filter expense accounts
	accounts, err := r.db.GetAccounts(userID)
	if err != nil {
		r.logger.Error("Failed to get accounts", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	expenseAccounts := make([]goserver.Account, 0)
	for _, a := range accounts {
		if a.Type == "expense" {
			expenseAccounts = append(expenseAccounts, a)
		}
	}
	data["ExpenseAccounts"] = expenseAccounts

	// Get planned items
	budgetItems, err := r.budgetService.ListMonthlyBudget(req.Context(), userID, monthStart)
	if err != nil {
		r.logger.Error("Failed to get budget items", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["BudgetItems"] = budgetItems

	// Create planned amounts map by account ID
	plannedByAcc := make(map[string]float64)
	for _, item := range budgetItems {
		plannedByAcc[item.AccountId] = item.Amount
	}
	data["PlannedByAccount"] = plannedByAcc

	// Get comparison (for actuals and totals)
	comparison, err := r.budgetService.CompareMonthly(req.Context(), userID, monthStart, "")
	if err != nil {
		r.logger.Error("Failed to get budget comparison", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Comparison"] = comparison

	// Build actual map from comparison rows
	actualByAcc := map[string]float64{}
	for _, row := range comparison.Rows {
		actualByAcc[row.AccountID] = row.Actual
	}
	data["ActualByAccount"] = actualByAcc

	// Execute the budget template
	if err := tmpl.ExecuteTemplate(w, "budget_unified.tpl", data); err != nil {
		r.logger.Error("budget_unified template exec failed", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (r *WebAppRouter) handleUnifiedBudgetPOST(w http.ResponseWriter, req *http.Request, userID string) {
	if err := req.ParseForm(); err != nil {
		r.logger.Error("Failed to parse form", "error", err)
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	// Parse month from form
	monthStart := normalizeMonth(time.Now())
	if monthParam := req.FormValue("month"); monthParam != "" {
		if ts, err := strconv.ParseInt(monthParam, 10, 64); err == nil {
			monthStart = time.Unix(ts, 0)
		}
	}
	monthStart = normalizeMonth(monthStart)

	if req.FormValue("action") == "copy" {
		fromMonth := monthStart.AddDate(0, -1, 0)
		count, err := r.budgetService.CopyFromPreviousMonth(req.Context(), userID, normalizeMonth(fromMonth), monthStart)
		if err != nil {
			r.logger.Error("Failed to copy budget", "error", err)
			redirectURL := "/web/budget?month=" + strconv.FormatInt(monthStart.Unix(), 10) + "&error=" + err.Error()
			http.Redirect(w, req, redirectURL, http.StatusSeeOther)
			return
		}
		successMsg := fmt.Sprintf("Copied %d budget items from %s", count, normalizeMonth(fromMonth).Format("January 2006"))
		redirectURL := "/web/budget?month=" + strconv.FormatInt(monthStart.Unix(), 10) + "&success=" + successMsg
		http.Redirect(w, req, redirectURL, http.StatusSeeOther)
		return
	}

	// Save plan
	var entries []goserver.BudgetItemNoId
	for key, values := range req.PostForm {
		if key == "month" || key == "action" {
			continue
		}
		if len(values) == 0 || values[0] == "" {
			continue
		}
		amount, err := strconv.ParseFloat(values[0], 64)
		if err != nil || amount <= 0 {
			continue
		}
		entries = append(entries, goserver.BudgetItemNoId{
			AccountId:   key,
			Amount:      amount,
			Date:        monthStart,
			Description: "Monthly budget",
		})
	}

	if err := r.budgetService.SaveMonthlyBudget(req.Context(), userID, monthStart, entries); err != nil {
		r.logger.Error("Failed to save budget", "error", err)
		redirectURL := "/web/budget?month=" + strconv.FormatInt(monthStart.Unix(), 10) + "&error=" + err.Error()
		http.Redirect(w, req, redirectURL, http.StatusSeeOther)
		return
	}

	redirectURL := "/web/budget?month=" + strconv.FormatInt(monthStart.Unix(), 10) + "&success=Budget saved successfully"
	http.Redirect(w, req, redirectURL, http.StatusSeeOther)
}
