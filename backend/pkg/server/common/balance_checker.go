package common

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func CheckBalanceForAccount(ctx context.Context, logger *slog.Logger, db database.Storage, userID, accountID string) error {
	logger.Info("Checking balance for account", "userID", userID, "accountID", accountID)

	acc, err := db.GetAccount(userID, accountID)
	if err != nil {
		return fmt.Errorf("failed to get account: %w", err)
	}

	count, err := db.CountUnprocessedTransactionsForAccount(userID, accountID, acc.IgnoreUnprocessedBefore)
	if err != nil {
		return fmt.Errorf("failed to count unprocessed transactions: %w", err)
	}

	if count > 0 {
		logger.Info("Account still has unprocessed transactions, skipping balance check", "count", count)
		return nil
	}

	for _, b := range acc.BankInfo.Balances {
		appBalance, err := db.GetAccountBalance(userID, accountID, b.CurrencyId)
		if err != nil {
			logger.With("error", err, "currencyId", b.CurrencyId).Error("Failed to get account balance")
			continue
		}

		if appBalance.Sub(b.ClosingBalance).Abs().GreaterThan(constants.ReconciliationTolerance) {
			logger.Warn("Balance mismatch detected",
				"account", acc.Name,
				"currencyId", b.CurrencyId,
				"appBalance", appBalance,
				"bankBalance", b.ClosingBalance)

			_, err := db.CreateNotification(userID, &goserver.Notification{
				Date:  time.Now(),
				Type:  string(models.NotificationTypeBalanceDoesntMatch),
				Title: "Balance Mismatch Detected",
				Description: fmt.Sprintf("Account %q has a balance mismatch. Account balance: %s, Bank balance: %s (Currency: %s). Please check your transactions.",
					acc.Name, appBalance.StringFixed(2), b.ClosingBalance.StringFixed(2), b.CurrencyId),
			})
			if err != nil {
				logger.With("error", err).Error("Failed to create balance mismatch notification")
			}
		} else {
			logger.Info("Balance verified for account", "account", acc.Name, "currencyId", b.CurrencyId, "balance", appBalance)
			// Create reconciliation record
			rec, err := db.CreateReconciliation(userID, &goserver.ReconciliationNoId{
				AccountId:         accountID,
				CurrencyId:        b.CurrencyId,
				ReconciledBalance: appBalance,
				ExpectedBalance:   b.ClosingBalance,
				IsManual:          false,
			})
			if err != nil {
				logger.With("error", err).Error("Failed to record reconciliation")
			} else {
				logger.Info("Reconciliation recorded", "recId", rec.ReconciliationId, "accountId", accountID)
			}
		}
	}

	return nil
}
