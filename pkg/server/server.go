package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/ya-breeze/geekbudgetbe/pkg/auth"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func Server(logger *slog.Logger, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	_, finishCham, err := Serve(ctx, logger, cfg)
	if err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	// Wait for an interrupt signal
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-stopChan
	logger.Info("Received signal. Shutting down server...")

	// Stop the server
	cancel()
	<-finishCham
	return nil
}

func createControllers(logger *slog.Logger, cfg *config.Config, db database.Storage) goserver.CustomControllers {
	return goserver.CustomControllers{
		AuthAPIService:         NewAuthAPIService(logger, db, cfg.JWTSecret),
		UserAPIService:         NewUserAPIService(logger, db),
		AccountsAPIService:     NewAccountsAPIService(logger, db),
		CurrenciesAPIService:   NewCurrenciesAPIServicer(logger, db),
		TransactionsAPIService: NewTransactionsAPIService(logger, db),
	}
}

func Serve(ctx context.Context, logger *slog.Logger, cfg *config.Config) (net.Addr, chan int, error) {
	commit := func() string {
		if info, ok := debug.ReadBuildInfo(); ok {
			for _, setting := range info.Settings {
				if setting.Key == "vcs.revision" {
					return setting.Value
				}
			}
		}

		return ""
	}()
	logger.Info("Built from git commit: " + commit)

	if cfg.JWTSecret == "" {
		logger.Warn("JWT secret is not set. Creating random secret...")
		cfg.JWTSecret = auth.GenerateRandomString(32)
	}

	storage := database.NewStorage(logger, cfg)
	if err := storage.Open(); err != nil {
		return nil, nil, fmt.Errorf("failed to open storage: %w", err)
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

	return goserver.Serve(ctx, logger, cfg,
		createControllers(logger, cfg, storage),
		[]goserver.Router{NewRootRouter(commit)},
		createMiddlewares(logger, cfg)...)
}

func upsertUser(storage database.Storage, username, hashedPassword string, logger *slog.Logger) error {
	user, err := storage.GetUser(username)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return fmt.Errorf("failed to reading user from DB: %w", err)
	}
	if user != nil {
		logger.Info(fmt.Sprintf("Updating password for user %q", username))
		user.HashedPassword = hashedPassword
		if err := storage.PutUser(user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		logger.Info(fmt.Sprintf("Creating user %q", username))
		if err := storage.CreateUser(username, hashedPassword); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	return nil
}

func createMiddlewares(logger *slog.Logger, cfg *config.Config) []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{
		AuthMiddleware(logger, cfg),
	}
}

type RootRouter struct {
	commit string
}

func NewRootRouter(commit string) *RootRouter {
	return &RootRouter{commit: commit}
}

func (r *RootRouter) Routes() goserver.Routes {
	return goserver.Routes{
		"RootPath": goserver.Route{
			Method:  "GET",
			Pattern: "/",
			HandlerFunc: func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintln(w, "GeekBudget API v0.0.1")
				fmt.Fprintln(w, "Git commit: "+r.commit)
			},
		},
	}
}
