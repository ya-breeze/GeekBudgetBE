package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

type budgetItemsAPIService struct {
	logger *slog.Logger
	db     database.Storage
}

func NewBudgetItemsAPIService(logger *slog.Logger, db database.Storage) goserver.BudgetItemsAPIServicer {
	return &budgetItemsAPIService{
		logger: logger,
		db:     db,
	}
}

// GetBudgetStatus - get budget status with rollover
func (s *budgetItemsAPIService) GetBudgetStatus(ctx context.Context, from time.Time, to time.Time) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	// Fetch all budget items (could optimize to fetch up to 'to' date)
	budgetItems, err := s.db.GetBudgetItems(userID)
	if err != nil {
		s.logger.Error("Failed to fetch budget items", "error", err)
		return goserver.Response(http.StatusInternalServerError, nil), err
	}

	// Calculate start date for rollover calculation (find earliest budget item)
	minDate := time.Now()
	if len(budgetItems) > 0 {
		minDate = budgetItems[0].Date
		for _, b := range budgetItems {
			if b.Date.Before(minDate) {
				minDate = b.Date
			}
		}
	}
	// Align minDate to start of month
	minDate = time.Date(minDate.Year(), minDate.Month(), 1, 0, 0, 0, 0, minDate.Location())

	// Fetch expenses from minDate to 'to'
	// We need all transactions to calculate actual spending
	// Note: Expenses logic might be complex (which accounts are expenses?).
	// Simplification: We look at transactions where money leaves a "budgeted" account (or matches category).
	// In this app, it seems BudgetItem is linked to AccountID.
	// So we check movements for that AccountID.
	// But usually BudgetItem AccountID refers to a "Category" account (Expense type).
	// Let's assume AccountID in BudgetItem is the Expense Account.

	transactions, err := s.db.GetTransactions(userID, minDate, to)
	if err != nil {
		s.logger.Error("Failed to fetch transactions", "error", err)
		return goserver.Response(http.StatusInternalServerError, nil), err
	}

	// Prepare data structures for calculation
	// Map: Month -> AccountID -> BudgetedAmount
	budgetMap := make(map[string]map[string]float64)
	// Map: Month -> AccountID -> SpentAmount
	spentMap := make(map[string]map[string]float64)

	// Helper to keys
	getMonthKey := func(d time.Time) string {
		return d.Format("2006-01")
	}

	for _, b := range budgetItems {
		key := getMonthKey(b.Date)
		if _, ok := budgetMap[key]; !ok {
			budgetMap[key] = make(map[string]float64)
		}
		// Sum in case multiple items for same month/account (though usually one)
		budgetMap[key][b.AccountId] += b.Amount
	}

	for _, t := range transactions {
		tMonth := getMonthKey(t.Date)
		for _, m := range t.Movements {
			// If movement is positive -> Income or refund?
			// If movement is negative -> Expense?
			// Usually expenses are negative movements on Expense accounts?
			// Or positive movements on Expense accounts (if Double Entry)?
			// Let's check typical usage.
			// In `prefillNewUser`: lunch is accGroceries (Amount 100), accCash (-100).
			// So expense account receives positive amount.
			if m.Amount > 0 {
				if _, ok := spentMap[tMonth]; !ok {
					spentMap[tMonth] = make(map[string]float64)
				}
				spentMap[tMonth][m.AccountId] += m.Amount
			}
		}
	}

	rolloverMap := make(map[string]float64) // AccountID -> Current Rollover

	// Iterate months from minDate to 'to'
	current := minDate

	results := []goserver.BudgetStatus{}

	// Ensure we cover up to 'to'
	// to date is exclusive, but we might want to include the month if it's partial?
	// The loop goes month by month.
	for current.Before(to) {
		monthKey := getMonthKey(current)

		// Get unique accounts involved this month (budgeted or spent)
		accounts := make(map[string]bool)
		for acc := range budgetMap[monthKey] {
			accounts[acc] = true
		}
		// We might want to track rollover for accounts even if no budget/spent this month?
		// Yes, if there is previous rollover.
		for acc := range rolloverMap {
			accounts[acc] = true
		}

		for accId := range accounts {
			budgeted := budgetMap[monthKey][accId]
			spent := spentMap[monthKey][accId]
			previousRollover := rolloverMap[accId]

			// Calculation
			// Available = Budget + Rollover
			// Remainder = Available - Spent
			// New Rollover = Remainder

			available := budgeted + previousRollover
			remainder := available - spent

			rolloverMap[accId] = remainder

			// Only add to results if within requested range
			if !current.Before(from) {
				results = append(results, goserver.BudgetStatus{
					Date:      current,
					AccountId: accId,
					Budgeted:  budgeted,
					Spent:     spent,
					Rollover:  previousRollover, // The rollover coming INTO this month
					Available: remainder,        // Remaining available for next month (or remaining now?)
					// Prompt: "visualize it comparing to the real expenses for this month"
					// Spec: Available
					// I'll return 'remainder' as available.
				})
			}
		}

		current = current.AddDate(0, 1, 0)
	}

	return goserver.Response(http.StatusOK, results), nil
}

// CreateBudgetItem - create new budgetItem
func (s *budgetItemsAPIService) CreateBudgetItem(ctx context.Context, budgetItemNoID goserver.BudgetItemNoId) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}
	budgetItem, err := s.db.CreateBudgetItem(userID, &budgetItemNoID)
	if err != nil {
		s.logger.Error("Failed to create budget item", "error", err)
		return goserver.Response(http.StatusInternalServerError, nil), err
	}
	return goserver.Response(http.StatusOK, budgetItem), nil
}

// DeleteBudgetItem - delete budgetItem
func (s *budgetItemsAPIService) DeleteBudgetItem(ctx context.Context, id string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}
	err := s.db.DeleteBudgetItem(userID, id)
	if err != nil {
		if err == database.ErrNotFound {
			return goserver.Response(http.StatusNotFound, nil), nil
		}
		s.logger.Error("Failed to delete budget item", "error", err)
		return goserver.Response(http.StatusInternalServerError, nil), err
	}
	return goserver.Response(http.StatusOK, nil), nil
}

// GetBudgetItem - get budgetItem
func (s *budgetItemsAPIService) GetBudgetItem(ctx context.Context, id string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}
	budgetItem, err := s.db.GetBudgetItem(userID, id)
	if err != nil {
		if err == database.ErrNotFound {
			return goserver.Response(http.StatusNotFound, nil), nil
		}
		s.logger.Error("Failed to get budget item", "error", err)
		return goserver.Response(http.StatusInternalServerError, nil), err
	}
	return goserver.Response(http.StatusOK, budgetItem), nil
}

// GetBudgetItems - get all budgetItems
func (s *budgetItemsAPIService) GetBudgetItems(ctx context.Context) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}
	budgetItems, err := s.db.GetBudgetItems(userID)
	if err != nil {
		s.logger.Error("Failed to get budget items", "error", err)
		return goserver.Response(http.StatusInternalServerError, nil), err
	}
	return goserver.Response(http.StatusOK, budgetItems), nil
}

// UpdateBudgetItem - update budgetItem
func (s *budgetItemsAPIService) UpdateBudgetItem(ctx context.Context, id string, budgetItemNoID goserver.BudgetItemNoId) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}
	budgetItem, err := s.db.UpdateBudgetItem(userID, id, &budgetItemNoID)
	if err != nil {
		if err == database.ErrNotFound {
			return goserver.Response(http.StatusNotFound, nil), nil
		}
		s.logger.Error("Failed to update budget item", "error", err)
		return goserver.Response(http.StatusInternalServerError, nil), err
	}
	return goserver.Response(http.StatusOK, budgetItem), nil
}
