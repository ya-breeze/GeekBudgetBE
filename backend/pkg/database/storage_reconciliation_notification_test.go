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

func TestReconciliationNotificationBehavior(t *testing.T) {
	// Setup
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := NewStorage(logger, cfg).(*storage)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-notify-1"
	currencyID := "CZK"

	// Create Currency
	cur, err := st.CreateCurrency(userID, &goserver.CurrencyNoId{Name: "CZK", Description: "CZK"})
	if err != nil {
		t.Fatalf("failed to create currency: %v", err)
	}
	currencyID = cur.Id

	// Create Account
	acc, err := st.CreateAccount(userID, &goserver.AccountNoId{
		Name:     "Notify Test Account",
		BankInfo: goserver.BankAccountInfo{Balances: []goserver.BankAccountInfoBalancesInner{{CurrencyId: currencyID, OpeningBalance: decimal.Zero}}},
	})
	if err != nil {
		t.Fatalf("failed to create account: %v", err)
	}
	accountID := acc.Id

	// Helper to create reconciliation manually
	createRec := func(reconciledAt time.Time) string {
		rec, err := st.CreateReconciliation(userID, &goserver.ReconciliationNoId{
			AccountId: accountID, CurrencyId: currencyID,
			ReconciledBalance: decimal.Zero, ExpectedBalance: decimal.Zero,
		})
		if err != nil {
			t.Fatalf("failed to create reconciliation: %v", err)
		}
		_ = st.db.Model(&models.Reconciliation{}).Where("id = ?", rec.ReconciliationId).Update("reconciled_at", reconciledAt).Error
		return rec.ReconciliationId
	}

	checkpoint := time.Now().Add(-24 * time.Hour)

	// Clear initial notifications if any
	_ = st.db.Where("user_id = ?", userID).Delete(&models.Notification{})

	t.Run("Create Transaction (Past) - No Notification", func(t *testing.T) {
		createRec(checkpoint)
		_, err := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        checkpoint.Add(-1 * time.Hour),
			Description: "Past Import",
			Movements:   []goserver.Movement{{AccountId: accountID, CurrencyId: currencyID, Amount: decimal.NewFromInt(10)}},
		})
		if err != nil {
			t.Fatalf("failed to create transaction: %v", err)
		}

		notifications, _ := st.GetNotifications(userID)
		if len(notifications) != 0 {
			t.Errorf("expected 0 notifications for import, got %d", len(notifications))
		}
	})

	t.Run("Update Transaction (User) - Notification Expected", func(t *testing.T) {
		// Create a transaction first (in the past)
		tr, err := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        checkpoint.Add(-2 * time.Hour),
			Description: "To be edited",
			Movements:   []goserver.Movement{{AccountId: accountID, CurrencyId: currencyID, Amount: decimal.NewFromInt(20)}},
		})
		if err != nil {
			t.Fatalf("failed to create initial transaction: %v", err)
		}

		createRec(checkpoint)
		// Clear notifications from initial transaction creation if any (should be 0 anyway)
		_ = st.db.Where("user_id = ?", userID).Delete(&models.Notification{})

		// Now update it (user style which preserves protected fields)

		// Now update it (user style which preserves protected fields)
		_, err = st.UpdateTransaction(userID, tr.Id, &goserver.TransactionNoId{
			Date:        checkpoint.Add(-2 * time.Hour),
			Description: "Edited by user",
			Movements:   []goserver.Movement{{AccountId: accountID, CurrencyId: currencyID, Amount: decimal.NewFromInt(30)}},
		})
		if err != nil {
			t.Fatalf("failed to update transaction: %v", err)
		}

		notifications, _ := st.GetNotifications(userID)
		if len(notifications) == 0 {
			t.Error("expected notification for user update, got none")
		}
		// Clean up for next subtest
		_ = st.db.Where("user_id = ?", userID).Delete(&models.Notification{})
	})

	t.Run("Update Transaction (Internal) - No Notification", func(t *testing.T) {
		tr, err := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        checkpoint.Add(-3 * time.Hour),
			Description: "Internal update target",
			Movements:   []goserver.Movement{{AccountId: accountID, CurrencyId: currencyID, Amount: decimal.NewFromInt(40)}},
		})
		if err != nil {
			t.Fatalf("failed to create initial transaction: %v", err)
		}

		createRec(checkpoint)
		_ = st.db.Where("user_id = ?", userID).Delete(&models.Notification{})

		// Update internal style

		// Update internal style
		_, err = st.UpdateTransactionInternal(userID, tr.Id, &goserver.TransactionNoId{
			Date:        checkpoint.Add(-3 * time.Hour),
			Description: "Internal update",
			Movements:   []goserver.Movement{{AccountId: accountID, CurrencyId: currencyID, Amount: decimal.NewFromInt(50)}},
		})
		if err != nil {
			t.Fatalf("failed to internal update transaction: %v", err)
		}

		notifications, _ := st.GetNotifications(userID)
		if len(notifications) != 0 {
			t.Errorf("expected 0 notifications for internal update, got %d", len(notifications))
		}
	})

	t.Run("Delete Transaction - Notification Expected", func(t *testing.T) {
		tr, err := st.CreateTransaction(userID, &goserver.TransactionNoId{
			Date:        checkpoint.Add(-4 * time.Hour),
			Description: "To be deleted",
			Movements:   []goserver.Movement{{AccountId: accountID, CurrencyId: currencyID, Amount: decimal.NewFromInt(60)}},
		})
		if err != nil {
			t.Fatalf("failed to create initial transaction: %v", err)
		}

		createRec(checkpoint)
		_ = st.db.Where("user_id = ?", userID).Delete(&models.Notification{})

		err = st.DeleteTransaction(userID, tr.Id)
		if err != nil {
			t.Fatalf("failed to delete transaction: %v", err)
		}

		notifications, _ := st.GetNotifications(userID)
		if len(notifications) == 0 {
			t.Error("expected notification for deletion, got none")
		}
	})
}
