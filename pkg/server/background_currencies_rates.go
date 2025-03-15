package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
)

func startCurrenciesRatesFetcher(
	ctx context.Context, logger *slog.Logger, db database.Storage,
) <-chan struct{} {
	logger.Info("Starting currencies rates fetcher...")

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ctx.Done():
				close(done)
				logger.Info("Stopped currencies rates fetcher")
				return
			default:
				// Do something
				logger.Info("Fetching currencies rates...")

				// importer := api.NewBankImportersAPIServiceImpl(logger, db)
				// pairs, err := db.GetAllBankImporters()
				// if err != nil {
				// 	logger.With("error", err).Error("Failed to get bank importers")

				// 	// Retry in 1 hour
				// 	select {
				// 	case <-time.After(time.Hour):
				// 		continue
				// 	case <-ctx.Done():
				// 		continue
				// 	}
				// }

				logger.Info("Delaying currencies rates fetcher for 24 hours...")
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
