package database

import (
	"log/slog"
	"testing"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestTransactionHistoryInternal(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := NewStorage(logger, cfg).(*storage)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-1"

	// 1. Create
	tr, err := st.CreateTransaction(userID, &goserver.TransactionNoId{
		Date:        time.Now(),
		Description: "History Test",
		Movements:   []goserver.Movement{},
	})
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	var history []models.AuditLog
	if err := st.db.Where("entity_id = ? AND entity_type = ?", tr.Id, "Transaction").Find(&history).Error; err != nil {
		t.Fatalf("failed to query audit log: %v", err)
	}

	if len(history) != 1 {
		t.Errorf("expected 1 audit record, got %d", len(history))
	} else if history[0].Action != "CREATED" {
		t.Errorf("expected CREATED action, got %s", history[0].Action)
	}

	// 2. Update
	_, err = st.UpdateTransaction(userID, tr.Id, &goserver.TransactionNoId{
		Date:        tr.Date,
		Description: "Updated History Test",
		Movements:   []goserver.Movement{},
	})
	if err != nil {
		t.Fatalf("failed to update transaction: %v", err)
	}

	if err := st.db.Where("entity_id = ? AND entity_type = ?", tr.Id, "Transaction").Find(&history).Error; err != nil {
		t.Fatalf("failed to query audit log: %v", err)
	}
	if len(history) != 2 {
		t.Errorf("expected 2 audit records, got %d", len(history))
	} else if history[1].Action != "UPDATED" {
		t.Errorf("expected UPDATED action, got %s", history[1].Action)
	}

	// 3. Delete
	err = st.DeleteTransaction(userID, tr.Id)
	if err != nil {
		t.Fatalf("failed to delete transaction: %v", err)
	}

	if err := st.db.Where("entity_id = ? AND entity_type = ?", tr.Id, "Transaction").Find(&history).Error; err != nil {
		t.Fatalf("failed to query audit log: %v", err)
	}
	if len(history) != 3 {
		t.Errorf("expected 3 audit records, got %d", len(history))
	} else if history[2].Action != "DELETED" {
		t.Errorf("expected DELETED action, got %s", history[2].Action)
	}
}
