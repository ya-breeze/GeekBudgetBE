package background

import (
	"context"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
)

type ForcedImportKeyType string

const ForcedImportKey ForcedImportKeyType = "forced_import_channel"

type ForcedImport struct {
	UserID         string
	BankImporterID string
}

func GetForcedImportChannel(ctx context.Context) chan<- ForcedImport {
	return ctx.Value(ForcedImportKey).(chan<- ForcedImport)
}

//nolint:gocognit,cyclop // TODO refactor
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
