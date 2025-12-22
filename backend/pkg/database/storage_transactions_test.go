package database_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestTransactionsStorage(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-1"

	t.Run("GetTransactions with onlySuspicious filter", func(t *testing.T) {
		t1 := goserver.Transaction{
			Date:              time.Now(),
			Description:       "Normal",
			SuspiciousReasons: []string{},
		}
		t2 := goserver.Transaction{
			Date:              time.Now(),
			Description:       "Suspicious",
			SuspiciousReasons: []string{"Too expensive"},
		}

		_, err := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              t1.Date,
			Description:       t1.Description,
			SuspiciousReasons: t1.SuspiciousReasons,
			Movements:         []goserver.Movement{},
		})
		if err != nil {
			t.Fatalf("failed to create t1: %v", err)
		}
		_, err = st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              t2.Date,
			Description:       t2.Description,
			SuspiciousReasons: t2.SuspiciousReasons,
			Movements:         []goserver.Movement{},
		})
		if err != nil {
			t.Fatalf("failed to create t2: %v", err)
		}

		// Test both=false
		all, err := st.GetTransactions(userID, time.Time{}, time.Time{}, false)
		if err != nil {
			t.Fatalf("failed to get all: %v", err)
		}
		if len(all) != 2 {
			t.Errorf("expected 2 transactions, got %d", len(all))
		}

		// Test onlySuspicious=true
		suspicious, err := st.GetTransactions(userID, time.Time{}, time.Time{}, true)
		if err != nil {
			t.Fatalf("failed to get suspicious: %v", err)
		}
		if len(suspicious) != 1 {
			t.Errorf("expected 1 suspicious transaction, got %d", len(suspicious))
		} else if suspicious[0].Description != "Suspicious" {
			t.Errorf("expected Suspicious, got %s", suspicious[0].Description)
		}
	})

	t.Run("TestSoftDeleteTransactions", func(t *testing.T) {
		// 1. Create a transaction
		date := time.Now()
		input := &goserver.TransactionNoId{
			Date:        date,
			Description: "Soft Delete Test",
			Movements:   []goserver.Movement{{Amount: 100, CurrencyId: "CZK", AccountId: "acc-1"}},
		}

		created, err := st.CreateTransaction(userID, input)
		if err != nil {
			t.Fatalf("failed to create transaction: %v", err)
		}

		// 2. Verify it exists
		all, err := st.GetTransactions(userID, time.Time{}, time.Time{}, false)
		if err != nil {
			t.Fatalf("failed to get transactions: %v", err)
		}
		found := false
		for _, tr := range all {
			if tr.Id == created.Id {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("transaction not found after creation")
		}

		// 3. Soft-delete it
		err = st.DeleteTransaction(userID, created.Id)
		if err != nil {
			t.Fatalf("failed to delete transaction: %v", err)
		}

		// 4. Verify GetTransactions does not return it
		allAfter, err := st.GetTransactions(userID, time.Time{}, time.Time{}, false)
		if err != nil {
			t.Fatalf("failed to get transactions after delete: %v", err)
		}
		for _, tr := range allAfter {
			if tr.Id == created.Id {
				t.Fatal("transaction still found in GetTransactions after soft-delete")
			}
		}

		// 5. Verify GetTransaction (singular) returns ErrNotFound
		_, err = st.GetTransaction(userID, created.Id)
		if err == nil {
			t.Fatal("GetTransaction should return error for soft-deleted record")
		}

		// 6. Verify GetTransactionsIncludingDeleted returns it
		allInc, err := st.GetTransactionsIncludingDeleted(userID, time.Time{}, time.Time{})
		if err != nil {
			t.Fatalf("failed to get transactions including deleted: %v", err)
		}
		found = false
		for _, tr := range allInc {
			if tr.Id == created.Id {
				found = true
				break
			}
		}
		if !found {
			t.Fatal("transaction not found in GetTransactionsIncludingDeleted after soft-delete")
		}
	})
}
