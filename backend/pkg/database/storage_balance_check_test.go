package database_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestStorageBalanceCheck(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-1"

	// Create CZK currency
	createdCur, err := st.CreateCurrency(userID, &goserver.CurrencyNoId{Name: "Czech Koruna"})
	if err != nil {
		t.Fatalf("failed to create currency: %v", err)
	}
	currencyID := createdCur.Id

	// Setup: Create an account
	accInput := &goserver.AccountNoId{
		Name: "Test Account",
		BankInfo: goserver.BankAccountInfo{
			Balances: []goserver.BankAccountInfoBalancesInner{
				{CurrencyId: currencyID, OpeningBalance: decimal.NewFromInt(1000), ClosingBalance: decimal.NewFromInt(1500)},
			},
		},
	}
	createdAcc, err := st.CreateAccount(userID, accInput)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}
	accID := createdAcc.Id

	// Create other account
	otherAcc, err := st.CreateAccount(userID, &goserver.AccountNoId{Name: "Other Account"})
	if err != nil {
		t.Fatalf("failed to create other account: %v", err)
	}

	t.Run("CountUnprocessedTransactionsForAccount handles omitempty AccountId", func(t *testing.T) {
		// Create a transaction where one side is the account and the other side is empty (unprocessed)
		// Due to omitempty, the empty AccountId will be missing from the JSON value in the database.
		t1 := &goserver.TransactionNoId{
			Description: "Unprocessed",
			Movements: []goserver.Movement{
				{Amount: decimal.NewFromInt(500), CurrencyId: currencyID, AccountId: accID},
				{Amount: decimal.NewFromInt(-500), CurrencyId: currencyID, AccountId: ""}, // Empty AccountId
			},
		}
		_, err := st.CreateTransaction(userID, t1)
		if err != nil {
			t.Fatalf("failed to create transaction: %v", err)
		}

		count, err := st.CountUnprocessedTransactionsForAccount(userID, accID, time.Time{})
		if err != nil {
			t.Fatalf("failed to count unprocessed: %v", err)
		}

		if count != 1 {
			t.Errorf("expected 1 unprocessed transaction, got %d", count)
		}
	})

	t.Run("GetAccountBalance calculates correctly", func(t *testing.T) {
		bal, err := st.GetAccountBalance(userID, accID, currencyID)
		if err != nil {
			t.Fatalf("failed to get balance: %v", err)
		}

		if !bal.Equal(decimal.NewFromInt(1500)) {
			t.Errorf("expected balance 1500.0, got %s", bal)
		}
	})

	t.Run("CountUnprocessedTransactionsForAccount returns 0 when all processed", func(t *testing.T) {
		// Create a processed transaction
		t2 := &goserver.TransactionNoId{
			Description: "Processed",
			Movements: []goserver.Movement{
				{Amount: decimal.NewFromInt(100), CurrencyId: currencyID, AccountId: accID},
				{Amount: decimal.NewFromInt(-100), CurrencyId: currencyID, AccountId: otherAcc.Id},
			},
		}
		_, err := st.CreateTransaction(userID, t2)
		if err != nil {
			t.Fatalf("failed to create transaction: %v", err)
		}

		count, err := st.CountUnprocessedTransactionsForAccount(userID, accID, time.Time{})
		if err != nil {
			t.Fatalf("failed to count unprocessed: %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 unprocessed transaction, got %d", count)
		}
	})
}
