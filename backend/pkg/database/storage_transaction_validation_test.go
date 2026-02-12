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

func TestTransactionValidation(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}
	st := database.NewStorage(logger, cfg)
	if err := st.Open(); err != nil {
		t.Fatalf("failed to open storage: %v", err)
	}
	defer st.Close()

	userID := "user-1"

	t.Run("Create transaction with non-existent account ID", func(t *testing.T) {
		input := &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Invalid Account",
			Movements: []goserver.Movement{
				{
					Amount:     decimal.NewFromInt(1),
					CurrencyId: "CZK",
					AccountId:  "non-existent-account",
				},
			},
		}

		_, err := st.CreateTransaction(userID, input)
		if err == nil {
			t.Error("expected error when creating transaction with non-existent account ID, but got nil")
		}
	})

	t.Run("Create transaction with non-existent currency ID", func(t *testing.T) {
		// First create a valid account to avoid account validation error
		acc, err := st.CreateAccount(userID, &goserver.AccountNoId{
			Name: "Test Account",
			Type: "cash",
		})
		if err != nil {
			t.Fatalf("failed to create account: %v", err)
		}

		input := &goserver.TransactionNoId{
			Date:        time.Now(),
			Description: "Invalid Currency",
			Movements: []goserver.Movement{
				{
					Amount:     decimal.NewFromInt(1),
					CurrencyId: "NON-EXISTENT",
					AccountId:  acc.Id,
				},
			},
		}

		_, err = st.CreateTransaction(userID, input)
		if err == nil {
			t.Error("expected error when creating transaction with non-existent currency ID, but got nil")
		}
	})
}
