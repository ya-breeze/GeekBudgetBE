package background

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

func StartDuplicateDetection(
	ctx context.Context, logger *slog.Logger, db database.Storage,
) <-chan struct{} {
	logger.Info("Starting duplicate detection task...")

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(done)
				logger.Info("Stopped duplicate detection task")
				return
			default:
				detectDuplicates(ctx, logger, db)

				logger.Info("Delaying duplicate detection for 24 hours...")
				select {
				case <-time.After(24 * time.Hour):
					continue
				case <-ctx.Done():
					continue
				}
			}
		}
	}()

	return done
}

func detectDuplicates(ctx context.Context, logger *slog.Logger, db database.Storage) {
	logger.Info("Running duplicate detection...")

	users, err := db.GetAllUserIDs()
	if err != nil {
		logger.With("error", err).Error("Failed to get users for duplicate detection")
		return
	}

	for _, userID := range users {
		processUserDuplicates(ctx, logger, db, userID)
	}

	logger.Info("Completed duplicate detection")
}

func processUserDuplicates(ctx context.Context, logger *slog.Logger, db database.Storage, userID string) {
	logger.Info("Processing duplicates for user", "userID", userID)

	// Fetch transactions from the last 30 days to avoid scanning everything
	dateFrom := time.Now().AddDate(0, 0, -30)
	transactions, err := db.GetTransactions(userID, dateFrom, time.Time{}, false)
	if err != nil {
		logger.With("error", err, "userID", userID).Error("Failed to get transactions for user")
		return
	}

	if len(transactions) < 2 {
		return
	}

	suspiciousCount := 0
	// Simple O(n^2) detection for small number of transactions in 30 days
	// For large number of transactions, we could optimize by sorting or grouping by amount
	for i := 0; i < len(transactions); i++ {
		t1 := transactions[i]
		if t1.DuplicateDismissed {
			continue
		}

		for j := i + 1; j < len(transactions); j++ {
			t2 := transactions[j]
			if t2.DuplicateDismissed {
				continue
			}

			// Check if they are duplicates
			if common.IsDuplicate(t1.Date, t1.Movements, t2.Date, t2.Movements) {
				// Check if they came from different sources (different external IDs)
				// If they have same external IDs, they were already handled by bank importer deduplication
				if hasDifferentSources(t1, t2) {
					// Link them together
					if err := db.AddDuplicateRelationship(userID, t1.Id, t2.Id); err != nil {
						logger.With("error", err, "t1", t1.Id, "t2", t2.Id).Error("Failed to add duplicate relationship")
					}

					marked := markAsSuspicious(logger, db, userID, &t1)
					if marked {
						suspiciousCount++
						// Update local copy to avoid double processing if it appears again
						t1.SuspiciousReasons = append(t1.SuspiciousReasons, models.DuplicateReason)
					}
					marked = markAsSuspicious(logger, db, userID, &t2)
					if marked {
						suspiciousCount++
						t2.SuspiciousReasons = append(t2.SuspiciousReasons, models.DuplicateReason)
					}
				}
			}
		}
	}

	if suspiciousCount > 0 {
		logger.Info("Detected duplicates for user", "userID", userID, "count", suspiciousCount)
		_, err := db.CreateNotification(userID, &goserver.Notification{
			Date:        time.Now(),
			Type:        string(models.NotificationTypeDuplicateDetected),
			Title:       "Potential Duplicate Transactions",
			Description: fmt.Sprintf("Detected %d potential duplicate transactions from different sources. Please review them.", suspiciousCount),
		})
		if err != nil {
			logger.With("error", err, "userID", userID).Error("Failed to create notification for duplicates")
		}
	}
}

func hasDifferentSources(t1, t2 goserver.Transaction) bool {
	if len(t1.ExternalIds) == 0 || len(t2.ExternalIds) == 0 {
		return false
	}

	// Check if they share any external ID
	for _, id1 := range t1.ExternalIds {
		for _, id2 := range t2.ExternalIds {
			if id1 == id2 {
				return false
			}
		}
	}
	return true
}

func markAsSuspicious(logger *slog.Logger, db database.Storage, userID string, t *goserver.Transaction) bool {
	if slices.Contains(t.SuspiciousReasons, models.DuplicateReason) {
		return false
	}

	t.SuspiciousReasons = append(t.SuspiciousReasons, models.DuplicateReason)
	tNoId := models.TransactionWithoutID(t)

	_, err := db.UpdateTransaction(userID, t.Id, tNoId)
	if err != nil {
		logger.With("error", err, "transactionID", t.Id).Error("Failed to update transaction with suspicious reason")
		return false
	}

	return true
}
