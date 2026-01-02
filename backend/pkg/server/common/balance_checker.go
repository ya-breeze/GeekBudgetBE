package common

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func CheckBalanceForAccount(ctx context.Context, logger *slog.Logger, db database.Storage, userID, accountID string) error {
	logger.Info("Checking balance for account", "userID", userID, "accountID", accountID)

	count, err := db.CountUnprocessedTransactionsForAccount(userID, accountID)
	if err != nil {
		return fmt.Errorf("failed to count unprocessed transactions: %w", err)
	}

	if count > 0 {
		logger.Info("Account still has unprocessed transactions, skipping balance check", "count", count)
		return nil
	}

	acc, err := db.GetAccount(userID, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	for _, b := range acc.BankInfo.Balances {
		appBalance, err := db.GetAccountBalance(userID, accountID, b.CurrencyId)
		if err != nil {
			logger.With("error", err, "currencyId", b.CurrencyId).Error("Failed to get app balance")
			continue
		}

		if math.Abs(appBalance-b.ClosingBalance) > 0.01 {
			logger.Warn("Balance mismatch detected",
				"account", acc.Name,
				"currencyId", b.CurrencyId,
				"appBalance", appBalance,
				"bankBalance", b.ClosingBalance)

			_, err := db.CreateNotification(userID, &goserver.Notification{
				Date:  time.Now(),
				Type:  string(models.NotificationTypeBalanceDoesntMatch),
				Title: "Balance Mismatch Detected",
				Description: fmt.Sprintf("Account %q has a balance mismatch. App balance: %.2f, Bank balance: %.2f (Currency: %s). Please check your transactions.",
					acc.Name, appBalance, b.ClosingBalance, b.CurrencyId),
			})
			if err != nil {
				logger.With("error", err).Error("Failed to create balance mismatch notification")
			}
		} else {
			logger.Info("Balance verified for account", "account", acc.Name, "currencyId", b.CurrencyId, "balance", appBalance)
		}
	}

	return nil
}
