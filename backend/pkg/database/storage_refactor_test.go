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

func TestStorageRefactor_DeleteAccount(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-1"

	// Create two accounts
	acc1, err := st.CreateAccount(userID, &goserver.AccountNoId{Name: "Account 1"})
	if err != nil {
		t.Fatalf("failed to create acc1: %v", err)
	}
	acc2, err := st.CreateAccount(userID, &goserver.AccountNoId{Name: "Account 2"})
	if err != nil {
		t.Fatalf("failed to create acc2: %v", err)
	}

	// Create currency
	curCZK, err := st.CreateCurrency(userID, &goserver.CurrencyNoId{Name: "Czech Koruna"})
	if err != nil {
		t.Fatalf("failed to create currency: %v", err)
	}

	// Create a transaction with a movement in acc1
	tr1Input := &goserver.TransactionNoId{
		Date:        time.Now(),
		Description: "Tr 1",
		Movements: []goserver.Movement{
			{Amount: decimal.NewFromInt(100), CurrencyId: curCZK.Id, AccountId: acc1.Id},
		},
	}
	tr1, err := st.CreateTransaction(userID, tr1Input)
	if err != nil {
		t.Fatalf("failed to create tr1: %v", err)
	}

	t.Run("Delete account without reassignment - should fail if in use", func(t *testing.T) {
		err := st.DeleteAccount(userID, acc1.Id, nil)
		if err != database.ErrAccountInUse {
			t.Errorf("expected ErrAccountInUse, got %v", err)
		}
	})

	t.Run("Delete account with reassignment", func(t *testing.T) {
		err := st.DeleteAccount(userID, acc1.Id, &acc2.Id)
		if err != nil {
			t.Fatalf("failed to delete account with reassignment: %v", err)
		}

		// Verify transaction now has acc2
		tr, err := st.GetTransaction(userID, tr1.Id)
		if err != nil {
			t.Fatalf("failed to get transaction: %v", err)
		}
		if tr.Movements[0].AccountId != acc2.Id {
			t.Errorf("expected account ID %s, got %s", acc2.Id, tr.Movements[0].AccountId)
		}

		// Verify account 1 is gone
		_, err = st.GetAccount(userID, acc1.Id)
		if err == nil {
			t.Error("account 1 should be deleted")
		}
	})
}

func TestStorageRefactor_DeleteCurrency(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-1"

	// Create two currencies
	cur1, err := st.CreateCurrency(userID, &goserver.CurrencyNoId{Name: "Currency 1"})
	if err != nil {
		t.Fatalf("failed to create cur1: %v", err)
	}
	cur2, err := st.CreateCurrency(userID, &goserver.CurrencyNoId{Name: "Currency 2"})
	if err != nil {
		t.Fatalf("failed to create cur2: %v", err)
	}

	// Create an account with balance in cur1 (JSON field bank_info)
	accInput := &goserver.AccountNoId{
		Name: "Account With Currency",
		BankInfo: goserver.BankAccountInfo{
			Balances: []goserver.BankAccountInfoBalancesInner{
				{OpeningBalance: decimal.NewFromInt(100), CurrencyId: cur1.Id},
			},
		},
	}
	acc, err := st.CreateAccount(userID, accInput)
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}

	// Create a transaction with movement in cur1 (JSON field movements)
	trInput := &goserver.TransactionNoId{
		Date:        time.Now(),
		Description: "Tr With Currency",
		Movements: []goserver.Movement{
			{Amount: decimal.NewFromInt(10), CurrencyId: cur1.Id, AccountId: acc.Id},
		},
	}
	tr, err := st.CreateTransaction(userID, trInput)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	t.Run("Delete currency without reassignment - should fail if in use", func(t *testing.T) {
		err := st.DeleteCurrency(userID, cur1.Id, nil)
		if err != database.ErrCurrencyInUse {
			t.Errorf("expected ErrCurrencyInUse, got %v", err)
		}
	})

	t.Run("Delete currency with reassignment", func(t *testing.T) {
		err := st.DeleteCurrency(userID, cur1.Id, &cur2.Id)
		if err != nil {
			t.Fatalf("failed to delete currency with reassignment: %v", err)
		}

		// Verify account now has cur2
		updatedAcc, err := st.GetAccount(userID, acc.Id)
		if err != nil {
			t.Fatalf("failed to get account: %v", err)
		}
		if updatedAcc.BankInfo.Balances[0].CurrencyId != cur2.Id {
			t.Errorf("expected currency ID %s in account, got %s", cur2.Id, updatedAcc.BankInfo.Balances[0].CurrencyId)
		}

		// Verify transaction now has cur2
		updatedTr, err := st.GetTransaction(userID, tr.Id)
		if err != nil {
			t.Fatalf("failed to get transaction: %v", err)
		}
		if updatedTr.Movements[0].CurrencyId != cur2.Id {
			t.Errorf("expected currency ID %s in transaction, got %s", cur2.Id, updatedTr.Movements[0].CurrencyId)
		}

		// Verify currency 1 is gone
		_, err = st.GetCurrency(userID, cur1.Id)
		if err == nil {
			t.Error("currency 1 should be deleted")
		}
	})
}
