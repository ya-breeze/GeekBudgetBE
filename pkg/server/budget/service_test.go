package budget

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/dusted-go/logging/prettylog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func createTestLogger() *slog.Logger {
	return slog.New(prettylog.NewHandler(&slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}))
}

func setupTestService(t *testing.T) (*Service, database.Storage) {
	logger := createTestLogger()
	cfg := &config.Config{}
	storage := database.NewStorage(logger, cfg)
	err := storage.Open()
	require.NoError(t, err)

	t.Cleanup(func() {
		storage.Close()
	})

	service := NewService(logger, storage)
	return service, storage
}

func TestValidateFutureMonth(t *testing.T) {
	service, _ := setupTestService(t)

	// Test: Current month should be rejected
	now := time.Now()
	currentMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	err := service.ValidateFutureMonth(currentMonth)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "budget month must be in the future")

	// Test: Past month should be rejected
	pastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	err = service.ValidateFutureMonth(pastMonth)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "budget month must be in the future")

	// Test: Next month should be accepted
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())
	err = service.ValidateFutureMonth(nextMonth)
	assert.NoError(t, err)

	// Test: Future month should be accepted
	futureMonth := time.Date(now.Year(), now.Month()+2, 1, 0, 0, 0, 0, now.Location())
	err = service.ValidateFutureMonth(futureMonth)
	assert.NoError(t, err)
}

func TestSaveMonthlyBudget_ValidationErrors(t *testing.T) {
	service, storage := setupTestService(t)
	ctx := context.Background()
	userID := "test-user"

	// Create test accounts
	expenseAccount := &goserver.AccountNoId{
		Name: "Food",
		Type: "expense",
	}
	incomeAccount := &goserver.AccountNoId{
		Name: "Salary",
		Type: "income",
	}

	expenseAcc, err := storage.CreateAccount(userID, expenseAccount)
	require.NoError(t, err)
	incomeAcc, err := storage.CreateAccount(userID, incomeAccount)
	require.NoError(t, err)

	// Past month is allowed (no error)
	now := time.Now()
	pastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())

	entries := []goserver.BudgetItemNoId{
		{
			AccountId:   expenseAcc.Id,
			Amount:      100.0,
			Date:        pastMonth,
			Description: "Test budget",
		},
	}

	err = service.SaveMonthlyBudget(ctx, userID, pastMonth, entries)
	assert.NoError(t, err)

	// Test: Non-expense account should be rejected
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	entriesWithIncome := []goserver.BudgetItemNoId{
		{
			AccountId:   incomeAcc.Id,
			Amount:      100.0,
			Date:        nextMonth,
			Description: "Test budget",
		},
	}

	err = service.SaveMonthlyBudget(ctx, userID, nextMonth, entriesWithIncome)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not an expense account")
}

func TestSaveMonthlyBudget_Success(t *testing.T) {
	service, storage := setupTestService(t)
	ctx := context.Background()
	userID := "test-user"

	// Create expense account
	expenseAccount := &goserver.AccountNoId{
		Name: "Food",
		Type: "expense",
	}

	expenseAcc, err := storage.CreateAccount(userID, expenseAccount)
	require.NoError(t, err)

	// Test: Valid future month with expense account should succeed
	now := time.Now()
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	entries := []goserver.BudgetItemNoId{
		{
			AccountId:   expenseAcc.Id,
			Amount:      100.0,
			Date:        nextMonth,
			Description: "Test budget",
		},
	}

	err = service.SaveMonthlyBudget(ctx, userID, nextMonth, entries)
	assert.NoError(t, err)

	// Verify budget was saved
	savedItems, err := service.ListMonthlyBudget(ctx, userID, nextMonth)
	require.NoError(t, err)
	assert.Len(t, savedItems, 1)
	assert.Equal(t, 100.0, savedItems[0].Amount)
	assert.Equal(t, expenseAcc.Id, savedItems[0].AccountId)
}

func TestCopyFromPreviousMonth_FillsOnlyZeroValues(t *testing.T) {
	service, storage := setupTestService(t)
	ctx := context.Background()
	userID := "test-user"

	// Create expense accounts
	expenseAcc1, err := storage.CreateAccount(userID, &goserver.AccountNoId{
		Name: "Groceries",
		Type: "expense",
	})
	require.NoError(t, err)

	expenseAcc2, err := storage.CreateAccount(userID, &goserver.AccountNoId{
		Name: "Transport",
		Type: "expense",
	})
	require.NoError(t, err)

	now := time.Now()
	fromMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, now.Location())
	toMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	// Set up source month with budget items
	fromEntries := []goserver.BudgetItemNoId{
		{
			AccountId:   expenseAcc1.Id,
			Amount:      300.0,
			Date:        fromMonth,
			Description: "Groceries budget",
		},
		{
			AccountId:   expenseAcc2.Id,
			Amount:      150.0,
			Date:        fromMonth,
			Description: "Transport budget",
		},
	}
	err = service.SaveMonthlyBudget(ctx, userID, fromMonth, fromEntries)
	require.NoError(t, err)

	// Set up target month with one non-zero item and one zero item
	toEntries := []goserver.BudgetItemNoId{
		{
			AccountId:   expenseAcc1.Id,
			Amount:      250.0, // Non-zero - should NOT be overwritten
			Date:        toMonth,
			Description: "Existing groceries budget",
		},
		{
			AccountId:   expenseAcc2.Id,
			Amount:      0.0, // Zero - should be filled from previous month
			Date:        toMonth,
			Description: "Empty transport budget",
		},
	}
	err = service.SaveMonthlyBudget(ctx, userID, toMonth, toEntries)
	require.NoError(t, err)

	// Copy from previous month
	count, err := service.CopyFromPreviousMonth(ctx, userID, fromMonth, toMonth)
	require.NoError(t, err)
	assert.Equal(t, 1, count) // Only transport should be copied (groceries has non-zero value)

	// Verify results - should have 3 items now (original 2 + 1 new from copy)
	resultItems, err := service.ListMonthlyBudget(ctx, userID, toMonth)
	require.NoError(t, err)
	require.Equal(t, 3, len(resultItems))

	// Group items by account ID
	groceriesItems := make([]goserver.BudgetItem, 0)
	transportItems := make([]goserver.BudgetItem, 0)
	for _, item := range resultItems {
		if item.AccountId == expenseAcc1.Id {
			groceriesItems = append(groceriesItems, item)
		} else if item.AccountId == expenseAcc2.Id {
			transportItems = append(transportItems, item)
		}
	}

	// Groceries should have only 1 item (non-zero value was not copied over)
	require.Equal(t, 1, len(groceriesItems))
	assert.Equal(t, 250.0, groceriesItems[0].Amount)
	assert.Equal(t, "Existing groceries budget", groceriesItems[0].Description)

	// Transport should have 2 items (original zero + new copied item)
	require.Equal(t, 2, len(transportItems))

	// Find the copied item (should have amount 150.0)
	var copiedTransportItem *goserver.BudgetItem
	for i := range transportItems {
		if transportItems[i].Amount == 150.0 {
			copiedTransportItem = &transportItems[i]
			break
		}
	}
	require.NotNil(t, copiedTransportItem)
	assert.Equal(t, 150.0, copiedTransportItem.Amount)
	assert.Equal(t, "Transport budget", copiedTransportItem.Description)
}
