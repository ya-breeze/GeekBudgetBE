package webapp

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

func (r *WebAppRouter) budgetComparisonHandler(w http.ResponseWriter, req *http.Request) {
	tmpl, err := r.loadTemplates()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data := utils.CreateTemplateData(req, "budget_comparison")

	userID, err := r.ValidateUserID(tmpl, w, req)
	if err != nil {
		r.logger.Error("Failed to get user ID from session", "error", err)
		return
	}
	data["UserID"] = userID

	// Parse month parameter (default to current month)
	monthStart := getCurrentMonth()
	if monthParam := req.URL.Query().Get("month"); monthParam != "" {
		if timestamp, err := strconv.ParseInt(monthParam, 10, 64); err == nil {
			monthStart = time.Unix(timestamp, 0)
		}
	}
	// Ensure we use first day of month
	monthStart = time.Date(monthStart.Year(), monthStart.Month(), 1, 0, 0, 0, 0, monthStart.Location())
	data["MonthStart"] = monthStart

	// Get budget comparison data
	comparison, err := r.budgetService.CompareMonthly(req.Context(), userID, monthStart, "")
	if err != nil {
		r.logger.Error("Failed to get budget comparison", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	data["Comparison"] = comparison

	// Add navigation data (URLs will be built in template)
	prevMonth := monthStart.AddDate(0, -1, 0)
	nextMonth := monthStart.AddDate(0, 1, 0)
	data["PrevMonth"] = prevMonth
	data["NextMonth"] = nextMonth

	if err := tmpl.ExecuteTemplate(w, "budget_comparison.tpl", data); err != nil {
		r.logger.Warn("failed to execute template", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getCurrentMonth returns the first day of current month
func getCurrentMonth() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
}
