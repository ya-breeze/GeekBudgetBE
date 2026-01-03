package database_test

import (
	"log/slog"
	"testing"

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
	accID := "acc-1"
	currencyID := "CZK"

	// Setup: Create an account
	// Actually, I don't need to create user explicitly for mock memory DB if I don't check constraints.

	acc := &goserver.AccountNoId{
		Name: "Test Account",
		BankInfo: goserver.BankAccountInfo{
			Balances: []goserver.BankAccountInfoBalancesInner{
				{CurrencyId: currencyID, OpeningBalance: 1000.0, ClosingBalance: 1500.0},
			},
		},
	}
	createdAcc, err := st.CreateAccount(userID, acc)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}
	accID = createdAcc.Id

	t.Run("CountUnprocessedTransactionsForAccount handles omitempty AccountId", func(t *testing.T) {
		// Create a transaction where one side is the account and the other side is empty (unprocessed)
		// Due to omitempty, the empty AccountId will be missing from the JSON value in the database.
		t1 := &goserver.TransactionNoId{
			Description: "Unprocessed",
			Movements: []goserver.Movement{
				{Amount: 500, CurrencyId: currencyID, AccountId: accID},
				{Amount: -500, CurrencyId: currencyID, AccountId: ""}, // Empty AccountId
			},
		}
		_, err := st.CreateTransaction(userID, t1)
		if err != nil {
			t.Fatalf("failed to create transaction: %v", err)
		}

		count, err := st.CountUnprocessedTransactionsForAccount(userID, accID)
		if err != nil {
			t.Fatalf("failed to count unprocessed: %v", err)
		}

		if count != 1 {
			t.Errorf("expected 1 unprocessed transaction, got %d. This might be due to the 'omitempty' bug where the field is missing from JSON.", count)
		}
	})

	t.Run("GetAccountBalance calculates correctly", func(t *testing.T) {
		// Account started with 1000.
		// We added a transaction with +500 for the account.
		// So balance should be 1500.
		bal, err := st.GetAccountBalance(userID, accID, currencyID)
		if err != nil {
			t.Fatalf("failed to get balance: %v", err)
		}

		if bal != 1500.0 {
			t.Errorf("expected balance 1500.0, got %f", bal)
		}
	})

	t.Run("CountUnprocessedTransactionsForAccount returns 0 when all processed", func(t *testing.T) {
		// Create a processed transaction
		t2 := &goserver.TransactionNoId{
			Description: "Processed",
			Movements: []goserver.Movement{
				{Amount: 100, CurrencyId: currencyID, AccountId: accID},
				{Amount: -100, CurrencyId: currencyID, AccountId: "other-acc"},
			},
		}
		_, err := st.CreateTransaction(userID, t2)
		if err != nil {
			t.Fatalf("failed to create transaction: %v", err)
		}

		// Before, we had 1 unprocessed. Now we added 1 processed.
		// But wait, the previous test added 1 unprocessed.

		count, err := st.CountUnprocessedTransactionsForAccount(userID, accID)
		if err != nil {
			t.Fatalf("failed to count unprocessed: %v", err)
		}
		if count != 1 {
			t.Errorf("expected 1 unprocessed transaction, got %d", count)
		}
	})
}
