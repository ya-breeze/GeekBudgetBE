package database

import (
	"log/slog"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestReconciliationHistoryPreservation(t *testing.T) {
	// Setup
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := NewStorage(logger, cfg).(*storage)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-fix-1"
	accountID := "acc-fix-1" // Placeholder, will be updated
	currencyID := "CZK"

	// Create Currency
	cur, err := st.CreateCurrency(userID, &goserver.CurrencyNoId{
		Name:        "Czech Koruna",
		Description: "CZK",
	})
	if err != nil {
		t.Fatalf("failed to create currency: %v", err)
	}
	currencyID = cur.Id

	// Create Account
	acc, err := st.CreateAccount(userID, &goserver.AccountNoId{
		Name:                    "Fix Test Account",
		Type:                    "CHECKING",
		BankInfo:                goserver.BankAccountInfo{Balances: []goserver.BankAccountInfoBalancesInner{{CurrencyId: currencyID, OpeningBalance: decimal.Zero}}},
		IgnoreUnprocessedBefore: time.Time{},
	})
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}
	accountID = acc.Id

	// Helper to create reconciliation manually (forcing ReconciledAt)
	createRec := func(reconciledAt time.Time, balance float64) string {
		recNoId := &goserver.ReconciliationNoId{
			AccountId:         accountID,
			CurrencyId:        currencyID,
			ReconciledBalance: decimal.NewFromFloat(balance),
			ExpectedBalance:   decimal.NewFromFloat(balance),
			IsManual:          true,
		}
		// We use CreateReconciliation but need to modify ReconciledAt manually in DB
		// because CreateReconciliation sets it to time.Now()
		rec, err := st.CreateReconciliation(userID, recNoId)
		if err != nil {
			t.Fatalf("failed to create reconciliation: %v", err)
		}

		err = st.db.Model(&models.Reconciliation{}).Where("id = ?", rec.ReconciliationId).Update("reconciled_at", reconciledAt).Error
		if err != nil {
			t.Fatalf("failed to update reconciled_at: %v", err)
		}
		return rec.ReconciliationId
	}

	// Timeline:
	// T1: Reconciliation 1 (Should be preserved)
	// T2: New Transaction inserted at T1.5 (Should trigger invalidation of T3)
	// T3: Reconciliation 2 (Should be deleted/invalidated)

	baseTime := time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	t1 := baseTime
	t1_5 := baseTime.Add(24 * time.Hour)
	t3 := baseTime.Add(48 * time.Hour)

	rec1ID := createRec(t1, 100.0)
	_ = createRec(t3, 200.0) // rec2ID

	// Verify both exist
	recs, err := st.GetReconciliationsForAccount(userID, accountID)
	if err != nil {
		t.Fatalf("failed to get reconciliations: %v", err)
	}
	if len(recs) != 2 {
		t.Fatalf("expected 2 reconciliations, got %d", len(recs))
	}

	// Insert Transaction at T1.5
	// This should invalidate reconciliations AFTER T1.5 (i.e. T3/rec2)
	// But PRESERVE T1/rec1
	_, err = st.CreateTransaction(userID, &goserver.TransactionNoId{
		Date:        t1_5,
		Description: "Retroactive Transaction",
		Movements: []goserver.Movement{
			{AccountId: accountID, CurrencyId: currencyID, Amount: decimal.NewFromFloat(50.0)},
		},
	})
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	// Verify State
	recsAfter, err := st.GetReconciliationsForAccount(userID, accountID)
	if err != nil {
		t.Fatalf("failed to get reconciliations after: %v", err)
	}

	// Log recs for debugging
	for _, r := range recsAfter {
		t.Logf("Remaining Rec: %s at %s", r.ReconciliationId, r.ReconciledAt)
	}

	if len(recsAfter) != 1 {
		t.Errorf("expected 1 reconciliation remaining, got %d", len(recsAfter))
		if len(recsAfter) == 0 {
			t.Error("ALL reconciliations were deleted!")
		}
	} else {
		if recsAfter[0].ReconciliationId != rec1ID {
			t.Errorf("expected Rec1 (%s) to remain, but got %s", rec1ID, recsAfter[0].ReconciliationId)
		}
	}
}
