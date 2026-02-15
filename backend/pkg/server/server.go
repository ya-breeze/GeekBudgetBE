package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/api"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/background"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/webapp"
	"github.com/ya-breeze/geekbudgetbe/pkg/version"
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

	if cfg.SessionSecret == "" {
		logger.Warn("Session secret is not set. Creating random secret...")
		cfg.SessionSecret = auth.GenerateRandomString(64)
	}

	logger.Info("Starting GeekBudget server...")

	if cfg.Users != "" {
		logger.Info("Creating users...")
		users := strings.Split(cfg.Users, ",")
		for _, user := range users {
			tokens := strings.Split(user, ":")
			if len(tokens) != 2 {
				return nil, nil, fmt.Errorf("invalid user format: %s", user)
			}

			if err := upsertUser(storage, tokens[0], tokens[1], logger); err != nil {
				return nil, nil, fmt.Errorf("failed to update user %q: %w", tokens[0], err)
			}
		}
	} else {
		logger.Info("No users defined in configuration")
	}

	controllers := createControllers(logger, cfg, storage)
	extraRouters := []goserver.Router{webapp.NewWebAppRouter(version.Commit, logger, cfg, storage)}
	extraRouters = append(extraRouters, api.NewCustomAuthAPIController(controllers.AuthAPIService, logger, cfg))
	extraRouters = append(extraRouters, api.NewStatusAPIController())

	return goserver.Serve(ctx, logger, cfg,
		controllers,
		extraRouters,
		createMiddlewares(logger, cfg, forcedImports)...)
}

func upsertUser(storage database.Storage, username, hashedPassword string, logger *slog.Logger) error {
	userID, err := storage.GetUserID(username)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return fmt.Errorf("failed to reading user from DB: %w", err)
	}
	var user *models.User
	if !errors.Is(err, database.ErrNotFound) {
		logger.Info(fmt.Sprintf("Updating password for user %q", username))

		user, err = storage.GetUser(userID)
		if err != nil {
			return fmt.Errorf("failed to get user: %w", err)
		}
		user.HashedPassword = hashedPassword
		if err = storage.PutUser(user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		logger.Info(fmt.Sprintf("Creating user %q", username))
		user, err = storage.CreateUser(username, hashedPassword)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	return nil
}

func createMiddlewares(
	logger *slog.Logger, cfg *config.Config, forcedImports chan<- common.ForcedImport,
) []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		CORSMiddleware(),
		AuthMiddleware(logger, cfg),
		ForcedImportMiddleware(logger, forcedImports),
	}
}
