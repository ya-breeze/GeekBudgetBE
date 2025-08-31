package webapp

import (
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
