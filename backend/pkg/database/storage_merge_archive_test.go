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

func TestMergeArchive(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-1"

	t.Run("Merge archives transaction data", func(t *testing.T) {
		// 1. Create two transactions
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Keep",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A1"}},
		})
		t2, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Merge",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A2"}},
		})

		// 2. Merge T2 into T1
		if _, err := st.MergeTransactions(userID, t1.Id, t2.Id); err != nil {
			t.Fatalf("failed to merge: %v", err)
		}

		// 3. Verify archive exists
		merged, err := st.GetMergedTransactions(userID)
		if err != nil {
			t.Fatalf("failed to get merged transactions: %v", err)
		}

		found := false
		for _, m := range merged {
			if m.Transaction.Id == t2.Id {
				found = true
				if m.MergedInto.Id != t1.Id {
					t.Errorf("expected merged into %s, got %s", t1.Id, m.MergedInto.Id)
				}
				if m.Transaction.Description != "Merge" {
					t.Errorf("expected description 'Merge', got %s", m.Transaction.Description)
				}
				break
			}
		}
		if !found {
			t.Error("merged transaction not found in archive")
		}

		// 4. Verify original transaction is hard-deleted (not just soft-deleted)
		_, err = st.GetTransaction(userID, t2.Id)
		if err == nil {
			t.Error("original transaction should be hard-deleted")
		}

		// 5. Verify GetTransaction for kept transaction shows merged ID
		updatedT1, _ := st.GetTransaction(userID, t1.Id)
		if len(updatedT1.MergedTransactionIds) != 1 || updatedT1.MergedTransactionIds[0] != t2.Id {
			t.Errorf("expected merged ID %s, got %v", t2.Id, updatedT1.MergedTransactionIds)
		}
	})

	t.Run("Delete duplicate archives transaction data", func(t *testing.T) {
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Keep",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A1"}},
		})
		t2, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Duplicate",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A2"}},
		})

		if err := st.DeleteDuplicateTransaction(userID, t2.Id, t1.Id); err != nil {
			t.Fatalf("failed to delete duplicate: %v", err)
		}

		merged, _ := st.GetMergedTransactions(userID)
		found := false
		for _, m := range merged {
			if m.Transaction.Id == t2.Id {
				found = true
				break
			}
		}
		if !found {
			t.Error("deleted duplicate not found in archive")
		}
	})

	t.Run("Unmerge cleans up archive", func(t *testing.T) {
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Keep",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A1"}},
		})
		t2, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Merge",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A2"}},
		})

		st.MergeTransactions(userID, t1.Id, t2.Id)

		if err := st.UnmergeTransaction(userID, t2.Id); err != nil {
			t.Fatalf("failed to unmerge: %v", err)
		}

		merged, _ := st.GetMergedTransactions(userID)
		for _, m := range merged {
			if m.Transaction.Id == t2.Id {
				t.Error("unmerged transaction still found in archive")
			}
		}

		// Verify restored
		restored, err := st.GetTransaction(userID, t2.Id)
		if err != nil {
			t.Fatalf("failed to get restored transaction: %v", err)
		}
		if restored.Id != t2.Id {
			t.Errorf("expected ID %s, got %s", t2.Id, restored.Id)
		}
	})

	t.Run("Unmerge preserves all transaction fields from archive", func(t *testing.T) {
		now := time.Now()
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        now,
			Description: "Keep",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A1"}},
			ExternalIds: []string{"ext1"},
		})
		t2, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              now,
			Description:       "Merge with fields",
			Place:             "Test Place",
			Tags:              []string{"tag1", "tag2"},
			PartnerName:       "Test Partner",
			PartnerAccount:    "12345",
			Extra:             "Extra data",
			Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A2"}},
			ExternalIds:       []string{"ext2"},
			SuspiciousReasons: []string{"test reason"},
		})

		st.MergeTransactions(userID, t1.Id, t2.Id)

		if err := st.UnmergeTransaction(userID, t2.Id); err != nil {
			t.Fatalf("failed to unmerge: %v", err)
		}

		// Verify all fields are preserved
		restored, err := st.GetTransaction(userID, t2.Id)
		if err != nil {
			t.Fatalf("failed to get restored transaction: %v", err)
		}

		if restored.Description != "Merge with fields" {
			t.Errorf("expected description 'Merge with fields', got %s", restored.Description)
		}
		if restored.Place != "Test Place" {
			t.Errorf("expected place 'Test Place', got %s", restored.Place)
		}
		if len(restored.Tags) != 2 {
			t.Errorf("expected 2 tags, got %d", len(restored.Tags))
		}
		if restored.PartnerName != "Test Partner" {
			t.Errorf("expected partner name 'Test Partner', got %s", restored.PartnerName)
		}
		if restored.PartnerAccount != "12345" {
			t.Errorf("expected partner account '12345', got %s", restored.PartnerAccount)
		}
		if restored.Extra != "Extra data" {
			t.Errorf("expected extra 'Extra data', got %s", restored.Extra)
		}
		if len(restored.ExternalIds) != 1 || restored.ExternalIds[0] != "ext2" {
			t.Errorf("expected external IDs [ext2], got %v", restored.ExternalIds)
		}
		if len(restored.SuspiciousReasons) != 1 || restored.SuspiciousReasons[0] != "test reason" {
			t.Errorf("expected suspicious reasons [test reason], got %v", restored.SuspiciousReasons)
		}
	})

	t.Run("Unmerge removes external IDs from kept transaction", func(t *testing.T) {
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Keep",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A1"}},
			ExternalIds: []string{"ext1", "ext_keep"},
		})
		t2, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Merge",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A2"}},
			ExternalIds: []string{"ext2"},
		})

		// Merge - this should transfer ext2 to t1
		st.MergeTransactions(userID, t1.Id, t2.Id)

		// Verify t1 has both external IDs
		kept, _ := st.GetTransaction(userID, t1.Id)
		if len(kept.ExternalIds) != 3 {
			t.Errorf("expected 3 external IDs after merge, got %d: %v", len(kept.ExternalIds), kept.ExternalIds)
		}

		// Unmerge - this should remove ext2 from t1
		if err := st.UnmergeTransaction(userID, t2.Id); err != nil {
			t.Fatalf("failed to unmerge: %v", err)
		}

		// Verify t1 only has original external IDs
		keptAfter, _ := st.GetTransaction(userID, t1.Id)
		if len(keptAfter.ExternalIds) != 2 {
			t.Errorf("expected 2 external IDs after unmerge, got %d: %v", len(keptAfter.ExternalIds), keptAfter.ExternalIds)
		}
		hasExt2 := false
		for _, id := range keptAfter.ExternalIds {
			if id == "ext2" {
				hasExt2 = true
			}
		}
		if hasExt2 {
			t.Error("ext2 should have been removed from kept transaction after unmerge")
		}
	})

	t.Run("Unmerge fails for non-merged transaction", func(t *testing.T) {
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Not merged",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A1"}},
		})

		err := st.UnmergeTransaction(userID, t1.Id)
		if err == nil {
			t.Error("expected error when unmerging non-merged transaction")
		}
	})

	t.Run("Unmerge fails for non-existent transaction", func(t *testing.T) {
		err := st.UnmergeTransaction(userID, "00000000-0000-0000-0000-000000000000")
		if err == nil {
			t.Error("expected error when unmerging non-existent transaction")
		}
	})

	t.Run("Merge and unmerge updates GetTransactions correctly", func(t *testing.T) {
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Keep for list test",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A1"}},
		})
		t2, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Merge for list test",
			Movements:   []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: "CZK", AccountId: "A2"}},
		})

		// Before merge - both should appear
		before, _ := st.GetTransactions(userID, time.Time{}, time.Time{}, false)
		countBefore := 0
		for _, tx := range before {
			if tx.Id == t1.Id || tx.Id == t2.Id {
				countBefore++
			}
		}
		if countBefore != 2 {
			t.Errorf("expected 2 transactions before merge, got %d", countBefore)
		}

		// After merge - only t1 should appear, with t2 in merged IDs
		st.MergeTransactions(userID, t1.Id, t2.Id)
		afterMerge, _ := st.GetTransactions(userID, time.Time{}, time.Time{}, false)
		foundT1 := false
		foundT2 := false
		for _, tx := range afterMerge {
			if tx.Id == t1.Id {
				foundT1 = true
				if len(tx.MergedTransactionIds) != 1 || tx.MergedTransactionIds[0] != t2.Id {
					t.Errorf("expected t1 to have t2 in merged IDs, got %v", tx.MergedTransactionIds)
				}
			}
			if tx.Id == t2.Id {
				foundT2 = true
			}
		}
		if !foundT1 {
			t.Error("t1 not found after merge")
		}
		if foundT2 {
			t.Error("t2 should not appear in GetTransactions after merge")
		}

		// After unmerge - both should appear again, t1 should have no merged IDs
		st.UnmergeTransaction(userID, t2.Id)
		afterUnmerge, _ := st.GetTransactions(userID, time.Time{}, time.Time{}, false)
		foundT1After := false
		foundT2After := false
		for _, tx := range afterUnmerge {
			if tx.Id == t1.Id {
				foundT1After = true
				if len(tx.MergedTransactionIds) != 0 {
					t.Errorf("expected t1 to have no merged IDs after unmerge, got %v", tx.MergedTransactionIds)
				}
			}
			if tx.Id == t2.Id {
				foundT2After = true
			}
		}
		if !foundT1After || !foundT2After {
			t.Error("both transactions should appear after unmerge")
		}
	})
}
