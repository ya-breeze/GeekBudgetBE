package test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func setupTestStorage(t *testing.T) database.Storage {
	logger := CreateTestLogger()
	cfg := &config.Config{}
	storage := database.NewStorage(logger, cfg)
	err := storage.Open()
	require.NoError(t, err)

	t.Cleanup(func() {
		storage.Close()
	})

	return storage
}

func TestGetBudgetItemsByMonth(t *testing.T) {
	storage := setupTestStorage(t)
	userID := "test-user"

	// Create test dates
	jan2024 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	feb2024 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	mar2024 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)

	// Create budget items across different months
	janItem := &goserver.BudgetItemNoId{
		AccountId:   "account1",
		Amount:      100.0,
		Date:        jan2024,
		Description: "January budget",
	}

	febItem := &goserver.BudgetItemNoId{
		AccountId:   "account1",
		Amount:      200.0,
		Date:        feb2024,
		Description: "February budget",
	}

	marItem := &goserver.BudgetItemNoId{
		AccountId:   "account2",
		Amount:      300.0,
		Date:        mar2024,
		Description: "March budget",
	}

	// Insert budget items
	_, err := storage.CreateBudgetItem(userID, janItem)
	require.NoError(t, err)

	_, err = storage.CreateBudgetItem(userID, febItem)
	require.NoError(t, err)

	_, err = storage.CreateBudgetItem(userID, marItem)
	require.NoError(t, err)

	// Test: Get January items only
	janItems, err := storage.GetBudgetItemsByMonth(userID, jan2024)
	require.NoError(t, err)
	assert.Len(t, janItems, 1)
	assert.Equal(t, "January budget", janItems[0].Description)
	assert.Equal(t, 100.0, janItems[0].Amount)

	// Test: Get February items only
	febItems, err := storage.GetBudgetItemsByMonth(userID, feb2024)
	require.NoError(t, err)
	assert.Len(t, febItems, 1)
	assert.Equal(t, "February budget", febItems[0].Description)
	assert.Equal(t, 200.0, febItems[0].Amount)

	// Test: Get March items only
	marItems, err := storage.GetBudgetItemsByMonth(userID, mar2024)
	require.NoError(t, err)
	assert.Len(t, marItems, 1)
	assert.Equal(t, "March budget", marItems[0].Description)
	assert.Equal(t, 300.0, marItems[0].Amount)

	// Test: Get items for month with no data
	apr2024 := time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC)
	aprItems, err := storage.GetBudgetItemsByMonth(userID, apr2024)
	require.NoError(t, err)
	assert.Len(t, aprItems, 0)
}

func TestCopyBudgetToMonth(t *testing.T) {
	storage := setupTestStorage(t)
	userID := "test-user"

	// Create test dates
	jan2024 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	feb2024 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	// Create budget items in January
	janItem1 := &goserver.BudgetItemNoId{
		AccountId:   "account1",
		Amount:      100.0,
		Date:        jan2024,
		Description: "January budget 1",
	}

	janItem2 := &goserver.BudgetItemNoId{
		AccountId:   "account2",
		Amount:      200.0,
		Date:        jan2024,
		Description: "January budget 2",
	}

	// Insert January budget items
	_, err := storage.CreateBudgetItem(userID, janItem1)
	require.NoError(t, err)

	_, err = storage.CreateBudgetItem(userID, janItem2)
	require.NoError(t, err)

	// Test: Copy January budget to February
	count, err := storage.CopyBudgetToMonth(userID, jan2024, feb2024)
	require.NoError(t, err)
	assert.Equal(t, 2, count)

	// Verify February items were created
	febItems, err := storage.GetBudgetItemsByMonth(userID, feb2024)
	require.NoError(t, err)
	assert.Len(t, febItems, 2)

	// Verify copied items have correct data
	for _, item := range febItems {
		assert.Equal(t, feb2024.Format("2006-01-02"), item.Date.Format("2006-01-02"))
		assert.True(t, item.Amount == 100.0 || item.Amount == 200.0)
		assert.True(t, item.AccountId == "account1" || item.AccountId == "account2")
	}

	// Verify original January items are unchanged
	janItems, err := storage.GetBudgetItemsByMonth(userID, jan2024)
	require.NoError(t, err)
	assert.Len(t, janItems, 2)
}
