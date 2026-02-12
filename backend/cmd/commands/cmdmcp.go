//nolint:forbidigo // it's okay to use fmt in this file
package commands

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/mcpserver"
)

func CmdMCP(log *slog.Logger) *cobra.Command {
	var username string

	res := &cobra.Command{
		Use:   "mcp",
		Short: "Run MCP stdio server for AI assistant integration",
		Long: `Start an MCP (Model Context Protocol) stdio server that allows AI assistants
like Claude Desktop or Claude Code to query GeekBudget data directly.

The server communicates via JSON-RPC over stdin/stdout and provides read-only
access to accounts, transactions, budgets, and other financial data.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Get config from context
			cfg, ok := cmd.Context().Value(ConfigKey).(*config.Config)
			if !ok {
				return errors.New("could not retrieve config from context")
			}

			// Create stderr-only logger (stdout is reserved for MCP JSON-RPC)
			logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.LevelInfo,
			}))

			// Initialize storage with stderr logger
			storage := database.NewStorage(logger, cfg)
			if err := storage.Open(); err != nil {
				return fmt.Errorf("failed to open storage: %w", err)
			}

			// Resolve username to userID
			userID, err := storage.GetUserID(username)
			if err != nil {
				return fmt.Errorf("failed to get user ID for username %q: %w", username, err)
			}

			logger.Info("Starting MCP server", "username", username, "userID", userID)

			// Run MCP server
			return mcpserver.Run(cmd.Context(), logger, storage, userID)
		},
		Args: cobra.NoArgs,
	}

	res.Flags().StringVarP(&username, "username", "u", "", "username for data access")
	_ = res.MarkFlagRequired("username")

	return res
}
