package database_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestDuplicateSynchronization(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-1"

	// Create CZK and EUR currencies
	curCZK, err := st.CreateCurrency(userID, &goserver.CurrencyNoId{Name: "Czech Koruna"})
	if err != nil {
		t.Fatalf("failed to create CZK: %v", err)
	}
	curEUR, err := st.CreateCurrency(userID, &goserver.CurrencyNoId{Name: "Euro"})
	if err != nil {
		t.Fatalf("failed to create EUR: %v", err)
	}

	// Create some accounts
	accA, _ := st.CreateAccount(userID, &goserver.AccountNoId{Name: "Acc A"})
	accB, _ := st.CreateAccount(userID, &goserver.AccountNoId{Name: "Acc B"})
	accC, _ := st.CreateAccount(userID, &goserver.AccountNoId{Name: "Acc C"})
	accR1, _ := st.CreateAccount(userID, &goserver.AccountNoId{Name: "Acc R1"})
	accR2, _ := st.CreateAccount(userID, &goserver.AccountNoId{Name: "Acc R2"})

	t.Run("Dismissal clears suspicious reason in linked transaction", func(t *testing.T) {
		// 1. Create two suspicious transactions
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              time.Now(),
			Description:       "T1",
			SuspiciousReasons: []string{models.DuplicateReason},
			Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: curCZK.Id, AccountId: accA.Id}},
		})
		t2, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              time.Now(),
			Description:       "T2",
			SuspiciousReasons: []string{models.DuplicateReason},
			Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: curCZK.Id, AccountId: accB.Id}},
		})

		// 2. Link them
		if err := st.AddDuplicateRelationship(userID, t1.Id, t2.Id); err != nil {
			t.Fatalf("failed to link: %v", err)
		}

		// 3. Dismiss T1
		update := models.TransactionWithoutID(&t1)
		update.DuplicateDismissed = true
		update.SuspiciousReasons = []string{} // Frontend does this
		if _, err := st.UpdateTransaction(userID, t1.Id, update); err != nil {
			t.Fatalf("failed to dismiss T1: %v", err)
		}

		// 4. Verify T2 no longer has the reason
		updatedT2, err := st.GetTransaction(userID, t2.Id)
		if err != nil {
			t.Fatalf("failed to get T2: %v", err)
		}
		for _, r := range updatedT2.SuspiciousReasons {
			if r == models.DuplicateReason {
				t.Errorf("T2 still has DuplicateReason after T1 dismissal")
			}
		}

		// 5. Verify links are gone
		links, _ := st.GetDuplicateTransactionIDs(userID, t2.Id)
		if len(links) > 0 {
			t.Errorf("T2 still has duplicate links: %v", links)
		}
	})

	t.Run("Merge clears suspicious reason in keep transaction", func(t *testing.T) {
		// 1. Create two suspicious transactions
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              time.Now(),
			Description:       "Keep",
			SuspiciousReasons: []string{models.DuplicateReason},
			ExternalIds:       []string{"ext1"},
			Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: curCZK.Id, AccountId: accA.Id}},
		})
		t2, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              time.Now(),
			Description:       "Merge",
			SuspiciousReasons: []string{models.DuplicateReason},
			ExternalIds:       []string{"ext2"},
			Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: curCZK.Id, AccountId: accB.Id}},
		})

		// 2. Link them
		st.AddDuplicateRelationship(userID, t1.Id, t2.Id)

		// 3. Merge T2 into T1
		if _, err := st.MergeTransactions(userID, t1.Id, t2.Id); err != nil {
			t.Fatalf("failed to merge: %v", err)
		}

		// 4. Verify T1 no longer has the reason and has merged IDs
		updatedT1, _ := st.GetTransaction(userID, t1.Id)
		for _, r := range updatedT1.SuspiciousReasons {
			if r == models.DuplicateReason {
				t.Errorf("T1 still has DuplicateReason after merge")
			}
		}
		if len(updatedT1.ExternalIds) != 2 {
			t.Errorf("expected 2 external IDs, got %v", updatedT1.ExternalIds)
		}

		// 5. Verify T2 is soft-deleted
		_, err := st.GetTransaction(userID, t2.Id)
		if err == nil {
			t.Errorf("T2 should be deleted after merge")
		}
	})

	t.Run("Reason stays if other duplicates exist", func(t *testing.T) {
		// T1 linked to T2 and T3
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              time.Now(),
			Description:       "T1",
			SuspiciousReasons: []string{models.DuplicateReason},
			Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: curCZK.Id, AccountId: accA.Id}},
		})
		t2, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              time.Now(),
			Description:       "T2",
			SuspiciousReasons: []string{models.DuplicateReason},
			Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: curCZK.Id, AccountId: accB.Id}},
		})
		t3, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              time.Now(),
			Description:       "T3",
			SuspiciousReasons: []string{models.DuplicateReason},
			Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(100), CurrencyId: curCZK.Id, AccountId: accC.Id}},
		})

		st.AddDuplicateRelationship(userID, t1.Id, t2.Id)
		st.AddDuplicateRelationship(userID, t1.Id, t3.Id)

		// Dismiss T2 -> Link T1-T2 gone, but T1-T3 stays.
		update2 := models.TransactionWithoutID(&t2)
		update2.DuplicateDismissed = true
		st.UpdateTransaction(userID, t2.Id, update2)

		// T1 should STILL have DuplicateReason because of T3
		updatedT1, _ := st.GetTransaction(userID, t1.Id)
		hasReason := false
		for _, r := range updatedT1.SuspiciousReasons {
			if r == models.DuplicateReason {
				hasReason = true
				break
			}
		}
		if !hasReason {
			t.Errorf("T1 lost DuplicateReason even though T3 link remains")
		}

		// Dismiss T3 -> Link T1-T3 gone.
		update3 := models.TransactionWithoutID(&t3)
		update3.DuplicateDismissed = true
		st.UpdateTransaction(userID, t3.Id, update3)

		// Now T1 should lose DuplicateReason
		updatedT1Final, _ := st.GetTransaction(userID, t1.Id)
		for _, r := range updatedT1Final.SuspiciousReasons {
			if r == models.DuplicateReason {
				t.Errorf("T1 still has DuplicateReason after all links gone")
			}
		}
	})

	t.Run("Revalidation clears link when date changes", func(t *testing.T) {
		now := time.Now()
		// Create two transactions with dates within 2 days
		t1, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              now,
			Description:       "Revalidate T1",
			SuspiciousReasons: []string{models.DuplicateReason},
			Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(500), CurrencyId: curEUR.Id, AccountId: accR1.Id}},
		})
		t2, _ := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:              now.Add(24 * time.Hour), // 1 day apart
			Description:       "Revalidate T2",
			SuspiciousReasons: []string{models.DuplicateReason},
			Movements:         []goserver.Movement{{Amount: decimal.NewFromInt(500), CurrencyId: curEUR.Id, AccountId: accR2.Id}},
		})

		// Link them
		st.AddDuplicateRelationship(userID, t1.Id, t2.Id)

		// Update T1's date to be 5 days earlier -> now > 2 days apart
		update := models.TransactionWithoutID(&t1)
		update.Date = now.Add(-5 * 24 * time.Hour)
		if _, err := st.UpdateTransaction(userID, t1.Id, update); err != nil {
			t.Fatalf("failed to update T1 date: %v", err)
		}

		// Verify link is gone
		links, _ := st.GetDuplicateTransactionIDs(userID, t1.Id)
		if len(links) > 0 {
			t.Errorf("T1 still has duplicate links after date change: %v", links)
		}

		// Verify both have DuplicateReason removed
		updatedT1, _ := st.GetTransaction(userID, t1.Id)
		for _, r := range updatedT1.SuspiciousReasons {
			if r == models.DuplicateReason {
				t.Errorf("T1 still has DuplicateReason after revalidation")
			}
		}
		updatedT2, _ := st.GetTransaction(userID, t2.Id)
		for _, r := range updatedT2.SuspiciousReasons {
			if r == models.DuplicateReason {
				t.Errorf("T2 still has DuplicateReason after revalidation")
			}
		}
	})
}
