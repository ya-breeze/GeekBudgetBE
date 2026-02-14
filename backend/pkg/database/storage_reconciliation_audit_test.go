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

func TestReconciliationAudit(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := NewStorage(logger, cfg).(*storage)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-1"
	accountID := "11111111-1111-1111-1111-111111111111"
	currencyID := "CZK"

	// 1. Create Reconciliation
	rec, err := st.CreateReconciliation(userID, &goserver.ReconciliationNoId{
		AccountId:         accountID,
		CurrencyId:        currencyID,
		ReconciledBalance: decimal.NewFromFloat(100.0),
		// ReconciledAt is not in ReconciliationNoId, it is set by storage
	})
	if err != nil {
		t.Fatalf("failed to create reconciliation: %v", err)
	}

	var history []models.AuditLog
	if err := st.db.Where("entity_id = ? AND entity_type = ?", rec.ReconciliationId, "Reconciliation").Find(&history).Error; err != nil {
		t.Fatalf("failed to query audit log: %v", err)
	}

	if len(history) != 1 {
		t.Errorf("expected 1 audit record, got %d", len(history))
	} else {
		if history[0].Action != "CREATED" {
			t.Errorf("expected CREATED action, got %s", history[0].Action)
		}
		if history[0].UserID != userID {
			t.Errorf("expected UserID %s, got %s", userID, history[0].UserID)
		}
	}

	// 2. Invalidate (Delete) Reconciliation
	if err := st.InvalidateReconciliation(userID, accountID, currencyID, time.Time{}); err != nil {
		t.Fatalf("failed to invalidate reconciliation: %v", err)
	}

	// Re-query history. We expect a new DELETED record.
	if err := st.db.Where("entity_id = ? AND entity_type = ?", rec.ReconciliationId, "Reconciliation").Order("created_at asc").Find(&history).Error; err != nil {
		t.Fatalf("failed to query audit log: %v", err)
	}

	if len(history) != 2 {
		t.Errorf("expected 2 audit records, got %d", len(history))
	} else {
		if history[1].Action != "DELETED" {
			t.Errorf("expected DELETED action, got %s", history[1].Action)
		}
	}
}
