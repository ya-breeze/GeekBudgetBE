package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/background"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/webapp"
	"github.com/ya-breeze/geekbudgetbe/pkg/version"
	kinauth "github.com/ya-breeze/kin-core/auth"
	"github.com/ya-breeze/kin-core/authdb"
	"gorm.io/gorm"
)

func Server(logger *slog.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	storage := database.NewStorage(logger, cfg)
	if err := storage.Open(); err != nil {
		return fmt.Errorf("failed to open storage: %w", err)
	}

	forcedImportChan := make(chan common.ForcedImport, 100)
	_, finishChan, err := Serve(ctx, logger, storage, cfg, forcedImportChan)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	// Start bank importers
	var importFinishChan <-chan struct{}
	if !cfg.DisableImporters {
		importFinishChan = background.StartBankImporters(ctx, logger, storage, cfg, forcedImportChan)
	} else {
		logger.Info("Bank importers are disabled")
	}

	// Start currency rate fetcher
	var fetchCurrenciesRatesChan <-chan struct{}
	if !cfg.DisableCurrenciesRatesFetch {
		fetchCurrenciesRatesChan = background.StartCurrenciesRatesFetcher(ctx, logger, storage)
	} else {
		logger.Info("Fetcher for currencies rates is disabled")
	}

	// Start duplicate detection
	duplicateDetectionFinishChan := background.StartDuplicateDetection(ctx, logger, storage)

	// Start database backup
	backupFinishChan := background.StartDatabaseBackup(ctx, logger, storage, cfg)

	// Wait for an interrupt signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan
	logger.Info("Received signal. Shutting down server...")

	// Stop the server
	cancel()
	<-finishChan
	if !cfg.DisableImporters {
		<-importFinishChan
	}
	if !cfg.DisableCurrenciesRatesFetch {
		<-fetchCurrenciesRatesChan
	}
	<-duplicateDetectionFinishChan
	<-backupFinishChan
	return nil
}

func createControllers(logger *slog.Logger, cfg *config.Config, db database.Storage) goserver.CustomControllers {
	unprocessedService := api.NewUnprocessedTransactionsAPIServiceImpl(logger, db)
	return goserver.CustomControllers{
		AuthAPIService:                    api.NewAuthAPIService(logger, db, cfg),
		UserAPIService:                    api.NewUserAPIService(logger, db),
		AccountsAPIService:                api.NewAccountsAPIService(logger, db, cfg),
		CurrenciesAPIService:              api.NewCurrenciesAPIServicer(logger, db),
		TransactionsAPIService:            api.NewTransactionsAPIService(logger, db),
		UnprocessedTransactionsAPIService: unprocessedService,
		MatchersAPIService:                api.NewMatchersAPIServiceImpl(logger, db, cfg, unprocessedService),
		BankImportersAPIService:           api.NewBankImportersAPIServiceImpl(logger, db, cfg),
		AggregationsAPIService:            api.NewAggregationsAPIServiceImpl(logger, db),
		NotificationsAPIService:           api.NewNotificationsAPIServiceImpl(logger, db),
		ImportAPIService:                  api.NewImportAPIServiceImpl(logger, db),
		ExportAPIService:                  api.NewExportAPIServiceImpl(logger, db),
		BudgetItemsAPIService:             api.NewBudgetItemsAPIService(logger, db),
		MergedTransactionsAPIService:      api.NewMergedTransactionsAPIService(logger, db),
		ReconciliationAPIService:          api.NewReconciliationAPIServiceImpl(logger, db),
		TemplatesAPIService:               api.NewTemplatesAPIServiceImpl(logger, db),
	}
}

func Serve(
	ctx context.Context, logger *slog.Logger,
	storage database.Storage, cfg *config.Config,
	forcedImports chan<- common.ForcedImport,
) (net.Addr, chan int, error) {
	// Initialize version info
	version.Init()
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				version.Commit = setting.Value
			case "vcs.time":
				version.BuildTime = setting.Value
			}
		}
	}
	logger.Info("Built from git commit: " + version.Commit)

	if cfg.JWTSecret == "" {
		logger.Warn("JWT secret is not set. Creating random secret...")
		cfg.JWTSecret = auth.GenerateRandomString(32)
	}

	logger.Info("Starting GeekBudget server...")

	gormDB := storage.GetDB()

	// Seed users from GB_SEED_USERS (format: "Family:Username:Password,...")
	if cfg.SeedUsers != "" {
		logger.Info("Seeding users...")
		for _, entry := range strings.Split(cfg.SeedUsers, ",") {
			entry = strings.TrimSpace(entry)
			if entry == "" {
				continue
			}
			if err := upsertSeedUser(storage, entry, logger); err != nil {
				return nil, nil, fmt.Errorf("failed to seed user %q: %w", entry, err)
			}
		}
	} else {
		logger.Info("No seed users defined in configuration")
	}

	// Start token cleanup goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				authdb.CleanupExpiredBlacklist(gormDB)
				authdb.CleanupExpiredRefreshTokens(gormDB)
			}
		}
	}()

	controllers := createControllers(logger, cfg, storage)
	extraRouters := []goserver.Router{webapp.NewWebAppRouter(version.Commit, logger, cfg, storage, gormDB)}
	extraRouters = append(extraRouters, api.NewCustomAuthAPIController(controllers.AuthAPIService, logger, cfg, storage, gormDB))
	extraRouters = append(extraRouters, api.NewStatusAPIController())

	return goserver.Serve(ctx, logger, cfg,
		controllers,
		extraRouters,
		createMiddlewares(logger, cfg, gormDB, forcedImports)...)
}

// upsertSeedUser creates or updates a user from a "Family:Username:Password" entry.
func upsertSeedUser(storage database.Storage, entry string, logger *slog.Logger) error {
	tokens := strings.Split(entry, ":")
	if len(tokens) != 3 {
		return fmt.Errorf("invalid seed user format %q, expected Family:Username:Password", entry)
	}
	familyName, username, password := tokens[0], tokens[1], tokens[2]

	// Ensure family exists
	family, err := storage.GetFamilyByName(familyName)
	if err != nil {
		family, err = storage.CreateFamily(familyName)
		if err != nil {
			return fmt.Errorf("failed to create family %q: %w", familyName, err)
		}
		logger.Info("Created family", "name", familyName, "id", family.ID)
	}

	// Upsert user
	existing, err := storage.GetUserByUsername(username)
	if err != nil {
		// User doesn't exist — create
		hash, hashErr := kinauth.HashPassword(password)
		if hashErr != nil {
			return fmt.Errorf("failed to hash password for %q: %w", username, hashErr)
		}
		user, createErr := storage.CreateUser(username, hash, family.ID)
		if createErr != nil {
			return fmt.Errorf("failed to create user %q: %w", username, createErr)
		}
		logger.Info("Created seed user", "username", username, "id", user.ID)
	} else {
		// User exists — update password
		hash, hashErr := kinauth.HashPassword(password)
		if hashErr != nil {
			return fmt.Errorf("failed to hash password for %q: %w", username, hashErr)
		}
		existing.PasswordHash = hash
		if putErr := storage.PutUser(existing); putErr != nil {
			return fmt.Errorf("failed to update user %q: %w", username, putErr)
		}
		logger.Info("Updated seed user password", "username", username)
	}

	return nil
}

func createMiddlewares(
	logger *slog.Logger, cfg *config.Config, gormDB *gorm.DB, forcedImports chan<- common.ForcedImport,
) []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		CORSMiddleware(),
		AuthMiddleware(logger, cfg, gormDB),
		ForcedImportMiddleware(logger, forcedImports),
	}
}
