package database_test

import (
	"log/slog"
	"testing"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func TestUpdateTransactionPreservesFields(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-update-test"

	// Create a transaction with rich data
	originalDate := time.Now().Add(-24 * time.Hour)
	input := &goserver.TransactionNoId{
		Date:               originalDate,
		Description:        "Original Description",
		ExternalIds:        []string{"ext-1", "ext-2"},
		UnprocessedSources: `{"some": "json"}`,
		SuspiciousReasons:  []string{"Duplicate candidate"},
		IsAuto:             true,
		Movements:          []goserver.Movement{},
	}

	created, err := st.CreateTransaction(userID, input)
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	// Verify created data
	if len(created.ExternalIds) != 2 {
		t.Errorf("expected 2 external IDs, got %d", len(created.ExternalIds))
	}

	// Update the transaction
	// - Simulate frontend sending only editable fields + DuplicateDismissed/IsAuto (if changed)
	// - We intentionally omit ExternalIds, UnprocessedSources, SuspiciousReasons
	updateInput := &goserver.TransactionNoId{
		Date:               originalDate, // Same date
		Description:        "Updated Description",
		DuplicateDismissed: true,  // Change this field
		IsAuto:             false, // Try to unset it (but it should be preserved!)
		Movements:          []goserver.Movement{},
	}

	updated, err := st.UpdateTransaction(userID, created.Id, updateInput)
	if err != nil {
		t.Fatalf("failed to update transaction: %v", err)
	}

	// Verify Description is updated
	if updated.Description != "Updated Description" {
		t.Errorf("expected Description to be 'Updated Description', got '%s'", updated.Description)
	}

	// Verify DuplicateDismissed is updated (Regression check)
	if !updated.DuplicateDismissed {
		t.Error("expected DuplicateDismissed to be updated to true")
	}

	// Verify IsAuto is PRESERVED (Safe method ignores input change)
	if !updated.IsAuto {
		t.Error("expected IsAuto to be preserved as true")
	}

	// Verify preserved fields
	if len(updated.ExternalIds) != 2 {
		t.Errorf("expected ExternalIds to be preserved (len=2), got %d", len(updated.ExternalIds))
	}
	if updated.UnprocessedSources != `{"some": "json"}` {
		t.Errorf("expected UnprocessedSources to be preserved, got '%s'", updated.UnprocessedSources)
	}
	if len(updated.SuspiciousReasons) != 1 {
		t.Errorf("expected SuspiciousReasons to be preserved (len=1), got %d", len(updated.SuspiciousReasons))
	}
}

func TestUpdateTransactionInternalUpdatesFields(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-internal-update-test"

	// Create a transaction with IsAuto=false
	originalDate := time.Now().Add(-24 * time.Hour)
	created, err := st.CreateTransaction(userID, &goserver.TransactionNoId{
		Date:        originalDate,
		Description: "Original Description",
		ExternalIds: []string{"ext-1"},
		IsAuto:      false,
		Movements:   []goserver.Movement{},
	})
	if err != nil {
		t.Fatalf("failed to create transaction: %v", err)
	}

	// Update with new values using Internal method
	updateInput := &goserver.TransactionNoId{
		Date:        originalDate,
		Description: "Updated Description",
		IsAuto:      true,              // This MUST be updated
		ExternalIds: []string{"ext-2"}, // This MUST also be updated (no preservation)
		Movements:   []goserver.Movement{},
	}

	updated, err := st.UpdateTransactionInternal(userID, created.Id, updateInput)
	if err != nil {
		t.Fatalf("failed to update transaction internal: %v", err)
	}

	// Verify IsAuto is updated (Unsafe method)
	if !updated.IsAuto {
		t.Error("expected IsAuto to be updated to true via Internal method")
	}
	// Verify ExternalID is updated (overwritten)
	if len(updated.ExternalIds) != 1 || updated.ExternalIds[0] != "ext-2" {
		t.Errorf("expected ExternalIds to be updated via Internal method, got %v", updated.ExternalIds)
	}
}
