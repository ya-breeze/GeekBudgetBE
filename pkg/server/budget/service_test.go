package budget

import (
	"context"
	"testing"
	"time"

	"log/slog"

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
