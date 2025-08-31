package budget

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestCompareMonthly(t *testing.T) {
	service, storage := setupTestService(t)
	ctx := context.Background()
	userID := "test-user"

	// Create test month
	testMonth := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)

	// Create expense accounts
	foodAccount := &goserver.AccountNoId{
		Name: "Food",
		Type: "expense",
	}
	transportAccount := &goserver.AccountNoId{
		Name: "Transport",
		Type: "expense",
	}

	foodAcc, err := storage.CreateAccount(userID, foodAccount)
	require.NoError(t, err)
	transportAcc, err := storage.CreateAccount(userID, transportAccount)
	require.NoError(t, err)

	// Create budget items (planned amounts)
	budgetItems := []goserver.BudgetItemNoId{
		{
			AccountId:   foodAcc.Id,
			Amount:      500.0, // Planned: $500 for food
			Date:        testMonth,
			Description: "Food budget",
		},
		{
			AccountId:   transportAcc.Id,
			Amount:      200.0, // Planned: $200 for transport
			Date:        testMonth,
			Description: "Transport budget",
		},
	}

	// Save budget items
	for _, item := range budgetItems {
		_, err := storage.CreateBudgetItem(userID, &item)
		require.NoError(t, err)
	}

	// Create actual transactions (actual expenses)
	// Food: $600 actual (over budget by $100)
	foodTransaction := &goserver.TransactionNoId{
		Date:        testMonth.AddDate(0, 0, 15), // Mid-month
		Description: "Grocery shopping",
		Movements: []goserver.Movement{
			{
				AccountId: foodAcc.Id,
				Amount:    600.0, // Actual: $600 spent on food
			},
		},
	}

	// Transport: $150 actual (under budget by $50)
	transportTransaction := &goserver.TransactionNoId{
		Date:        testMonth.AddDate(0, 0, 20),
		Description: "Bus pass",
		Movements: []goserver.Movement{
			{
				AccountId: transportAcc.Id,
				Amount:    150.0, // Actual: $150 spent on transport
			},
		},
	}

	// Save transactions
	_, err = storage.CreateTransaction(userID, foodTransaction)
	require.NoError(t, err)
	_, err = storage.CreateTransaction(userID, transportTransaction)
	require.NoError(t, err)

	// Test: Compare monthly budget vs actual
	comparison, err := service.CompareMonthly(ctx, userID, testMonth, "")
	require.NoError(t, err)

	// Verify totals
	assert.Equal(t, 700.0, comparison.TotalPlanned) // $500 + $200
	assert.Equal(t, 750.0, comparison.TotalActual)  // $600 + $150
	assert.Equal(t, 50.0, comparison.TotalDelta)    // $750 - $700 = $50 over budget

	// Verify individual account rows
	assert.Len(t, comparison.Rows, 2)

	// Find food and transport rows
	var foodRow, transportRow *Row
	for i := range comparison.Rows {
		if comparison.Rows[i].AccountName == "Food" {
			foodRow = &comparison.Rows[i]
		} else if comparison.Rows[i].AccountName == "Transport" {
			transportRow = &comparison.Rows[i]
		}
	}

	// Verify food account
	require.NotNil(t, foodRow)
	assert.Equal(t, foodAcc.Id, foodRow.AccountID)
	assert.Equal(t, "Food", foodRow.AccountName)
	assert.Equal(t, 500.0, foodRow.Planned)
	assert.Equal(t, 600.0, foodRow.Actual)
	assert.Equal(t, 100.0, foodRow.Delta) // Over budget by $100

	// Verify transport account
	require.NotNil(t, transportRow)
	assert.Equal(t, transportAcc.Id, transportRow.AccountID)
	assert.Equal(t, "Transport", transportRow.AccountName)
	assert.Equal(t, 200.0, transportRow.Planned)
	assert.Equal(t, 150.0, transportRow.Actual)
	assert.Equal(t, -50.0, transportRow.Delta) // Under budget by $50
}

func TestCompareMonthly_NoBudgetItems(t *testing.T) {
	service, storage := setupTestService(t)
	ctx := context.Background()
	userID := "test-user"

	// Create test month with no budget items
	testMonth := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)

	// Create expense account
	foodAccount := &goserver.AccountNoId{
		Name: "Food",
		Type: "expense",
	}

	foodAcc, err := storage.CreateAccount(userID, foodAccount)
	require.NoError(t, err)

	// Create actual transaction (no budget planned)
	transaction := &goserver.TransactionNoId{
		Date:        testMonth.AddDate(0, 0, 15),
		Description: "Unplanned expense",
		Movements: []goserver.Movement{
			{
				AccountId: foodAcc.Id,
				Amount:    100.0,
			},
		},
	}

	_, err = storage.CreateTransaction(userID, transaction)
	require.NoError(t, err)

	// Test: Compare with no budget items
	comparison, err := service.CompareMonthly(ctx, userID, testMonth, "")
	require.NoError(t, err)

	// Verify totals
	assert.Equal(t, 0.0, comparison.TotalPlanned)
	assert.Equal(t, 100.0, comparison.TotalActual)
	assert.Equal(t, 100.0, comparison.TotalDelta) // All spending is over budget

	// Verify account row shows actual spending
	assert.Len(t, comparison.Rows, 1)
	assert.Equal(t, "Food", comparison.Rows[0].AccountName)
	assert.Equal(t, 0.0, comparison.Rows[0].Planned)
	assert.Equal(t, 100.0, comparison.Rows[0].Actual)
	assert.Equal(t, 100.0, comparison.Rows[0].Delta)
}
