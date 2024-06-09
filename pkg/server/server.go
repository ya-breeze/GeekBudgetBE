package server

import (
	"fmt"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func Server(logger *slog.Logger, cfg *config.Config) error {
	storage := database.NewStorage(logger, cfg)
	if err := storage.Open(); err != nil {
		return fmt.Errorf("Failed to open storage: %w", err)
	}

	logger.Info(fmt.Sprintf("Starting GeekBudget at port %d...", cfg.Port))
	return goserver.Serve(cfg)

	// userId := "123e4567-e89b-12d3-a456-426614174000"

	// acc := &goserver.AccountNoId{
	// 	Name:        "Test Account",
	// 	Type:        "CHECKING",
	// 	Description: "Test Account Description",
	// }

	// account, err := storage.CreateAccount(userId, acc)
	// if err != nil {
	// 	return fmt.Errorf("Failed to create account: %w", err)
	// }

	// logger.With("account", account).Info("Account created")

	// accounts, err := storage.GetAccounts(userId)
	// if err != nil {
	// 	return fmt.Errorf("Failed to get accounts: %w", err)
	// }

	// logger.With("accounts", accounts).Info(fmt.Sprintf("Accounts retrieved: %d", len(accounts)))

	// logger.Info("GeekBudget stopped")

	// return nil
}
