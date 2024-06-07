package main

import (
	"log/slog"
	"os"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Starting GeekBudget...")

	cfg := &config.Config{}

	storage := database.NewStorage(logger, cfg)
	if storage.Open() != nil {
		logger.Error("Failed to open storage")
		return
	}

	userId := "123e4567-e89b-12d3-a456-426614174000"

	acc := &goserver.AccountNoId{
		Name:        "Test Account",
		Type:        "CHECKING",
		Description: "Test Account Description",
	}

	account, err := storage.CreateAccount(userId, acc)
	if err != nil {
		logger.Error("Failed to create account")
		return
	}

	logger.With("account", account).Info("Account created")

	accounts, err := storage.GetAccounts(userId)
	if err != nil {
		logger.Error("Failed to get accounts")
		return
	}

	logger.With("accounts", accounts).Info("Accounts retrieved")

	logger.Info("GeekBudget stopped")
}
