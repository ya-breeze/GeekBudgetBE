package api

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/shopspring/decimal"
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
func (s *budgetItemsAPIService) GetBudgetStatus(ctx context.Context, from time.Time, to time.Time, outputCurrencyId string, includeHidden bool) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(http.StatusInternalServerError, nil), nil
	}

	// Fetch all budget items
	budgetItems, err := s.db.GetBudgetItems(userID)
	if err != nil {
		s.logger.Error("Failed to fetch budget items", "error", err)
		return goserver.Response(http.StatusInternalServerError, nil), err
	}

	// Fetch all accounts to determine their primary currency
	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.Error("Failed to fetch accounts", "error", err)
		return goserver.Response(http.StatusInternalServerError, nil), err
	}
	accountCurrencyMap := make(map[string]string) // AccountID -> CurrencyID
	allowedAccounts := make(map[string]bool)
	for _, acc := range accounts {
		if includeHidden || !acc.HideFromReports {
			allowedAccounts[acc.Id] = true
		}
		// Use first balance currency as "Account Currency" for budgeting purposes
		if len(acc.BankInfo.Balances) > 0 {
			accountCurrencyMap[acc.Id] = acc.BankInfo.Balances[0].CurrencyId
		}
	}

	// Helpers for currency conversion
	currencyMap := buildCurrencyMap(s.logger, s.db, userID)
	outputCurrencyName := ""
	if outputCurrencyId != "" {
		outputCurrencyName = currencyMap[outputCurrencyId]
	}
	currenciesRatesFetcher := common.NewCurrenciesRatesFetcher(s.logger, s.db)

	// Calculate start date for rollover calculation (find earliest budget item)
	minDate := time.Now()
	if len(budgetItems) > 0 {
		minDate = budgetItems[0].Date
		for _, b := range budgetItems {
			if !allowedAccounts[b.AccountId] {
				continue
			}
			if b.Date.Before(minDate) {
				minDate = b.Date
			}
		}
	}
	// Align minDate to start of month
	minDate = time.Date(minDate.Year(), minDate.Month(), 1, 0, 0, 0, 0, minDate.Location())

	transactions, err := s.db.GetTransactions(userID, minDate, to, false)
	if err != nil {
		s.logger.Error("Failed to fetch transactions", "error", err)
		return goserver.Response(http.StatusInternalServerError, nil), err
	}

	// Map: Month -> AccountID -> BudgetedAmount (Converted)
	budgetMap := make(map[string]map[string]decimal.Decimal)
	// Map: Month -> AccountID -> SpentAmount (Converted)
	spentMap := make(map[string]map[string]decimal.Decimal)

	// Helper to keys
	getMonthKey := func(d time.Time) string {
		return d.Format("2006-01")
	}

	for _, b := range budgetItems {
		if !allowedAccounts[b.AccountId] {
			continue
		}
		key := getMonthKey(b.Date)
		if _, ok := budgetMap[key]; !ok {
			budgetMap[key] = make(map[string]decimal.Decimal)
		}

		amount := b.Amount
		// Convert if needed
		accCurrencyId := accountCurrencyMap[b.AccountId]
		if outputCurrencyId != "" && accCurrencyId != "" && accCurrencyId != outputCurrencyId {
			// Find currency name
			originalCurrencyName := currencyMap[accCurrencyId]
			converted, err := currenciesRatesFetcher.Convert(ctx, b.Date, originalCurrencyName, outputCurrencyName, amount)
			if err == nil {
				amount = converted
			} else {
				s.logger.Warn("Failed to convert budget amount", "error", err, "from", originalCurrencyName, "to", outputCurrencyName)
			}
		}

		budgetMap[key][b.AccountId] = budgetMap[key][b.AccountId].Add(amount)
	}

	for _, t := range transactions {
		tMonth := getMonthKey(t.Date)
		for _, m := range t.Movements {
			if !allowedAccounts[m.AccountId] {
				continue
			}
			if m.Amount.IsPositive() {
				if _, ok := spentMap[tMonth]; !ok {
					spentMap[tMonth] = make(map[string]decimal.Decimal)
				}

				amount := m.Amount
				if outputCurrencyId != "" && m.CurrencyId != outputCurrencyId {
					originalCurrencyName := currencyMap[m.CurrencyId]
					converted, err := currenciesRatesFetcher.Convert(ctx, t.Date, originalCurrencyName, outputCurrencyName, amount)
					if err == nil {
						amount = converted
					} else {
						s.logger.Warn("Failed to convert movement amount", "error", err)
					}
				}

				spentMap[tMonth][m.AccountId] = spentMap[tMonth][m.AccountId].Add(amount)
			}
		}
	}

	rolloverMap := make(map[string]decimal.Decimal) // AccountID -> Current Rollover

	// Iterate months from minDate to 'to'
	current := minDate
	results := []goserver.BudgetStatus{}

	for current.Before(to) {
		monthKey := getMonthKey(current)

		// Get unique accounts involved this month (budgeted or spent)
		accountsSet := make(map[string]bool)
		for acc := range budgetMap[monthKey] {
			accountsSet[acc] = true
		}
		for acc := range spentMap[monthKey] {
			accountsSet[acc] = true
		}
		for acc := range rolloverMap {
			accountsSet[acc] = true
		}

		for accId := range accountsSet {
			budgeted := budgetMap[monthKey][accId]
			spent := spentMap[monthKey][accId]
			previousRollover := rolloverMap[accId]

			available := budgeted.Add(previousRollover)
			remainder := available.Sub(spent)

			rolloverMap[accId] = remainder

			// Only add to results if within requested range
			if !current.Before(from) {
				results = append(results, goserver.BudgetStatus{
					Date:      current,
					AccountId: accId,
					Budgeted:  budgeted,
					Spent:     spent,
					Rollover:  previousRollover,
					Available: remainder,
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
