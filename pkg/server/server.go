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
	"time"

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
		AuthAPIService:                    NewAuthAPIService(logger, db, cfg.JWTSecret),
		UserAPIService:                    NewUserAPIService(logger, db),
		AccountsAPIService:                NewAccountsAPIService(logger, db),
		CurrenciesAPIService:              NewCurrenciesAPIServicer(logger, db),
		TransactionsAPIService:            NewTransactionsAPIService(logger, db),
		UnprocessedTransactionsAPIService: NewUnprocessedTransactionsAPIServiceImpl(logger, db),
		MatchersAPIService:                NewMatchersAPIServiceImpl(logger, db),
		BankImportersAPIService:           NewBankImportersAPIServiceImpl(logger, db),
		AggregationsAPIService:            NewAggregationsAPIServiceImpl(logger, db),
		NotificationsAPIService:           NewNotificationsAPIServiceImpl(logger, db),
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
		if err = storage.PutUser(user); err != nil {
			return fmt.Errorf("failed to update user: %w", err)
		}
	} else {
		logger.Info(fmt.Sprintf("Creating user %q", username))
		if user, err = storage.CreateUser(username, hashedPassword); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
	}

	err = prefillNewUser(storage, user.ID.String())
	if err != nil {
		return fmt.Errorf("failed to prefill new user: %w", err)
	}

	return nil
}

//nolint:funlen // This function is long because it creates many default items
func prefillNewUser(storage database.Storage, userID string) error {
	// Create default accounts
	account := &goserver.AccountNoId{
		Name:        "Cash",
		Description: "Cash account",
		Type:        "asset",
	}
	accCash, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create cash account: %w", err)
	}

	account = &goserver.AccountNoId{
		Name:        "Bank",
		Description: "Bank account",
		Type:        "asset",
	}
	accBank, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create bank account: %w", err)
	}

	account = &goserver.AccountNoId{
		Name:        "Salary",
		Description: "Salary account",
		Type:        "income",
	}
	accSalary, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create income account: %w", err)
	}

	account = &goserver.AccountNoId{
		Name: "🥩 Food",
		Type: "expense",
	}
	accFood, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create food account: %w", err)
	}

	account = &goserver.AccountNoId{
		Name:        "Transport",
		Description: "Transport expenses account",
		Type:        "expense",
	}
	accTransport, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create transport account: %w", err)
	}

	account = &goserver.AccountNoId{
		Name: "🏠 Rent",
		Type: "expense",
	}
	accRent, err := storage.CreateAccount(userID, account)
	if err != nil {
		return fmt.Errorf("failed to create rent account: %w", err)
	}

	// Create default currencies
	currency := &goserver.CurrencyNoId{
		Name: "USD",
	}
	curUSD, err := storage.CreateCurrency(userID, currency)
	if err != nil {
		return fmt.Errorf("failed to create USD currency: %w", err)
	}

	currency = &goserver.CurrencyNoId{
		Name: "EUR",
	}
	curEUR, err := storage.CreateCurrency(userID, currency)
	if err != nil {
		return fmt.Errorf("failed to create EUR currency: %w", err)
	}

	currency = &goserver.CurrencyNoId{
		Name: "CZK",
	}
	curCZK, err := storage.CreateCurrency(userID, currency)
	if err != nil {
		return fmt.Errorf("failed to create CZK currency: %w", err)
	}

	// Create default transactions
	transaction := &goserver.TransactionNoId{
		Date:        time.Now().Add(-5 * 24 * time.Hour),
		Description: "Initial state for cash",
		Tags:        []string{"initial_account_state"},
		Movements: []goserver.Movement{
			{
				AccountId:  accCash.Id,
				Amount:     1000,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accCash.Id,
				Amount:     100,
				CurrencyId: curUSD.Id,
			},
			{
				AccountId:  accCash.Id,
				Amount:     10,
				CurrencyId: curEUR.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create initial cash transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-5 * 24 * time.Hour),
		Description: "Initial state for bank",
		Tags:        []string{"initial_account_state"},
		Movements: []goserver.Movement{
			{
				AccountId:  accBank.Id,
				Amount:     10000,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create initial bank transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-24 * time.Hour),
		Description: "Monthly salary",
		Movements: []goserver.Movement{
			{
				AccountId:  accSalary.Id,
				Amount:     -10000,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accBank.Id,
				Amount:     10000,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create salary transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-4 * time.Hour),
		Description: "Lunch",
		Movements: []goserver.Movement{
			{
				AccountId:  accFood.Id,
				Amount:     100,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accCash.Id,
				Amount:     -100,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create lunch transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-24 * time.Hour),
		Description: "Lunch",
		Movements: []goserver.Movement{
			{
				AccountId:  accFood.Id,
				Amount:     150,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accCash.Id,
				Amount:     -150,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create lunch transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-3 * time.Hour),
		Description: "Bus ticket",
		Movements: []goserver.Movement{
			{
				AccountId:  accTransport.Id,
				Amount:     25,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accCash.Id,
				Amount:     -25,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create bus ticket transaction: %w", err)
	}

	transaction = &goserver.TransactionNoId{
		Date:        time.Now().Add(-2 * time.Hour),
		Description: "Apartment rent",
		Movements: []goserver.Movement{
			{
				AccountId:  accRent.Id,
				Amount:     5000,
				CurrencyId: curCZK.Id,
			},
			{
				AccountId:  accBank.Id,
				Amount:     -5000,
				CurrencyId: curCZK.Id,
			},
		},
	}
	if _, err := storage.CreateTransaction(userID, transaction); err != nil {
		return fmt.Errorf("failed to create rent transaction: %w", err)
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
