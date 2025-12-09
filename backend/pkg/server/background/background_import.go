package background

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
)

type ForcedImportKeyType string

const ForcedImportKey ForcedImportKeyType = "forced_import_channel"

type ForcedImport struct {
	UserID         string
	BankImporterID string
}

func GetForcedImportChannel(ctx context.Context) chan<- ForcedImport {
	res, _ := ctx.Value(ForcedImportKey).(chan<- ForcedImport)
	return res
}

//nolint:funlen,gocognit,cyclop // TODO refactor
func StartBankImporters(
	ctx context.Context, logger *slog.Logger, db database.Storage, forcedImports <-chan ForcedImport,
) <-chan struct{} {
	logger.Info("Starting bank importers...")

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(done)
				logger.Info("Stopped bank importers")
				return
			case forcedImport := <-forcedImports:
				logger.Info("Forced import", "userID", forcedImport.UserID, "BankImporterID", forcedImport.BankImporterID)
			default:
				// Do something
				logger.Info("Importing from bank importers...")

				importer := api.NewBankImportersAPIServiceImpl(logger, db)
				pairs, err := db.GetAllBankImporters()
				if err != nil {
					logger.With("error", err).Error("Failed to get bank importers")

					// Retry in 1 hour
					select {
					case forcedImport := <-forcedImports:
						logger.Info("Forced import", "userID", forcedImport.UserID, "BankImporterID", forcedImport.BankImporterID)
					case <-time.After(time.Hour):
						continue
					case <-ctx.Done():
						continue
					}
				}

				for _, pair := range pairs {
					if pair.BankImporterType != "fio" {
						logger.Info("Skipping bank importer type", "type", pair.BankImporterType)
						continue
					}
					i, err := importer.Fetch(ctx, pair.UserID, pair.BankImporterID)
					if err != nil {
						logger.With("error", err).Error("Failed to import bank transactions")
						continue
					}

					if i != nil {
						logger.Info("Imported bank transactions successfully", "result", i)
					}
				}

				// Process unprocessed transactions for auto-conversion before delay
				processUnprocessedTransactionsForAutoConversion(ctx, logger, db)

				logger.Info("Delaying bank imports for 24 hours...")
				select {
				case forcedImport := <-forcedImports:
					logger.Info("Forced import", "userID", forcedImport.UserID, "BankImporterID", forcedImport.BankImporterID)
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

// processUnprocessedTransactionsForAutoConversion processes all unprocessed transactions
// and automatically converts those that have exactly one matcher with 100% success history
func processUnprocessedTransactionsForAutoConversion(
	ctx context.Context, logger *slog.Logger, db database.Storage,
) {
	logger.Info("Processing unprocessed transactions for auto-conversion...")

	// Get all users to process their unprocessed transactions
	users, err := getAllUsers(db)
	if err != nil {
		logger.With("error", err).Error("Failed to get users for auto-conversion")
		return
	}

	unprocessedService := api.NewUnprocessedTransactionsAPIServiceImpl(logger, db)

	for _, userID := range users {
		logger.Info("Processing unprocessed transactions for user", "userID", userID)

		// Get all unprocessed transactions for this user
		unprocessedTransactions, _, err := unprocessedService.PrepareUnprocessedTransactions(
			ctx, userID, false, "",
		)
		if err != nil {
			logger.With("error", err, "userID", userID).Error("Failed to get unprocessed transactions")
			continue
		}

		if len(unprocessedTransactions) == 0 {
			logger.Info("No unprocessed transactions for user", "userID", userID)
			continue
		}

		logger.Info("Found unprocessed transactions for auto-conversion",
			"userID", userID, "count", len(unprocessedTransactions))

		// Process each unprocessed transaction
		for _, unprocessed := range unprocessedTransactions {
			processUnprocessedTransactionForAutoConversion(
				ctx, logger, db, unprocessedService, userID, unprocessed,
			)
		}
	}

	logger.Info("Completed processing unprocessed transactions for auto-conversion")
}

// processUnprocessedTransactionForAutoConversion processes a single unprocessed transaction
// for potential auto-conversion based on matcher success history
func processUnprocessedTransactionForAutoConversion(
	ctx context.Context, logger *slog.Logger, db database.Storage,
	unprocessedService *api.UnprocessedTransactionsAPIServiceImpl,
	userID string, unprocessed goserver.UnprocessedTransaction,
) {
	if len(unprocessed.Matched) == 0 {
		logger.Info("No matched matchers for transaction",
			"transactionID", unprocessed.Transaction.Id, "userID", userID)
		return
	}

	perfectMatchers := findPerfectMatchers(db, logger, userID, unprocessed.Matched)

	switch len(perfectMatchers) {
	case 1:
		matcher := perfectMatchers[0]
		logger.Info("Auto-converting unprocessed transaction using perfect matcher",
			"transactionID", unprocessed.Transaction.Id,
			"matcherID", matcher.MatcherId,
			"userID", userID)

		// Convert the transaction using the perfect matcher
		convertedTransaction, err := unprocessedService.Convert(
			ctx, userID, unprocessed.Transaction.Id, &matcher.Transaction,
		)
		if err != nil {
			logger.With("error", err, "transactionID", unprocessed.Transaction.Id,
				"matcherID", matcher.MatcherId, "userID", userID).Error(
				"Failed to auto-convert unprocessed transaction")
			return
		}

		// Add successful confirmation to the matcher's history
		if err := db.AddMatcherConfirmation(userID, matcher.MatcherId, true); err != nil {
			logger.With("error", err, "matcherID", matcher.MatcherId, "userID", userID).Warn(
				"Failed to add confirmation to matcher after auto-conversion")
		}

		logger.Info("Successfully auto-converted unprocessed transaction",
			"transactionID", convertedTransaction.Id,
			"matcherID", matcher.MatcherId,
			"userID", userID)
	case 0:
		logger.Debug("No matchers with 100% success history for transaction",
			"transactionID", unprocessed.Transaction.Id,
			"userID", userID)
	default:
		logger.Debug("Multiple matchers with 100% success history, keeping transaction unprocessed",
			"transactionID", unprocessed.Transaction.Id,
			"userID", userID,
			"perfectMatchersCount", len(perfectMatchers))
	}
}

// findPerfectMatchers returns matchers (from matchedList) whose confirmation
// history exists and contains only successful confirmations (all true).
func findPerfectMatchers(
	db database.Storage, logger *slog.Logger, userID string,
	matchedList []goserver.MatcherAndTransaction,
) []goserver.MatcherAndTransaction {
	perfect := make([]goserver.MatcherAndTransaction, 0, len(matchedList))

	for _, matched := range matchedList {
		matcher, err := db.GetMatcher(userID, matched.MatcherId)
		if err != nil {
			logger.With("error", err, "matcherID", matched.MatcherId, "userID", userID).Warn(
				"Failed to get matcher for auto-conversion check")
			continue
		}

		history := matcher.GetConfirmationHistory()
		if len(history) == 0 {
			logger.Info("Matcher has no confirmation history",
				"matcherID", matched.MatcherId, "userID", userID)
			continue
		}

		allSuccessful := true
		for _, confirmed := range history {
			if !confirmed {
				allSuccessful = false
				break
			}
		}

		if len(history) < 10 {
			logger.Info("Matcher has insufficient confirmation history",
				"matcherID", matched.MatcherId, "userID", userID,
				"historyLength", len(history))
			continue
		}

		if allSuccessful {
			perfect = append(perfect, matched)
			logger.Info("Found matcher with 100% success history",
				"matcherID", matched.MatcherId, "userID", userID,
				"historyLength", len(history))
		}
	}

	return perfect
}

// ProcessUnprocessedTransactionsForAutoConversion is an exported wrapper used by
// tests and external callers to trigger the auto-conversion pass.
func ProcessUnprocessedTransactionsForAutoConversion(
	ctx context.Context, logger *slog.Logger, db database.Storage,
) {
	processUnprocessedTransactionsForAutoConversion(ctx, logger, db)
}

// getAllUsers retrieves all user IDs from the database
func getAllUsers(db database.Storage) ([]string, error) {
	// Get all bank importers and extract unique user IDs
	// This is a simple way to get active users - could be optimized with a dedicated method
	importers, err := db.GetAllBankImporters()
	if err != nil {
		return nil, fmt.Errorf("failed to get bank importers to extract users: %w", err)
	}

	userSet := make(map[string]bool)
	for _, importer := range importers {
		userSet[importer.UserID] = true
	}

	users := make([]string, 0, len(userSet))
	for userID := range userSet {
		users = append(users, userID)
	}

	return users, nil
}
