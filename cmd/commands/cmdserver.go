package commands

import (
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

func CmdServer(logger *slog.Logger, cfg *config.Config) *cobra.Command {
	res := &cobra.Command{
		Use:   "server",
		Short: "Start HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return server(logger, cfg)
		},
	}

	return res
}

func server(logger *slog.Logger, cfg *config.Config) error {
	storage := database.NewStorage(logger, cfg)
	if err := storage.Open(); err != nil {
		return fmt.Errorf("Failed to open storage: %w", err)
	}

	logger.Info("Starting GeekBudget...")

	userId := "123e4567-e89b-12d3-a456-426614174000"

	acc := &goserver.AccountNoId{
		Name:        "Test Account",
		Type:        "CHECKING",
		Description: "Test Account Description",
	}

	account, err := storage.CreateAccount(userId, acc)
	if err != nil {
		return fmt.Errorf("Failed to create account: %w", err)
	}

	logger.With("account", account).Info("Account created")

	accounts, err := storage.GetAccounts(userId)
	if err != nil {
		return fmt.Errorf("Failed to get accounts: %w", err)
	}

	logger.With("accounts", accounts).Info(fmt.Sprintf("Accounts retrieved: %d", len(accounts)))

	logger.Info("GeekBudget stopped")

	return nil
}
