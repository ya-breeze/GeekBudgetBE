//nolint:forbidigo // it's okay to use fmt in this file
package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
)

func CmdMCPConfig(log *slog.Logger) *cobra.Command {
	var username, output string

	res := &cobra.Command{
		Use:   "mcp-config",
		Short: "Generate or update .mcp.json configuration file",
		Long: `Generate an MCP configuration file that can be used by AI assistants
like Claude Desktop or Claude Code to connect to GeekBudget.

The generated configuration includes the absolute path to the geekbudget binary
and the database location. If the output file already exists, it will be updated
while preserving other MCP server configurations.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			// Get config from context
			cfg, ok := cmd.Context().Value(ConfigKey).(*config.Config)
			if !ok {
				return errors.New("could not retrieve config from context")
			}

			// Get absolute path to current executable
			execPath, err := os.Executable()
			if err != nil {
				return fmt.Errorf("failed to get executable path: %w", err)
			}
			execPath, err = filepath.EvalSymlinks(execPath)
			if err != nil {
				return fmt.Errorf("failed to resolve executable path: %w", err)
			}

			// Get absolute database path
			dbPath := cfg.DBPath
			if !filepath.IsAbs(dbPath) {
				dbPath, err = filepath.Abs(dbPath)
				if err != nil {
					return fmt.Errorf("failed to resolve database path: %w", err)
				}
			}

			// Read existing config if it exists
			conf := make(map[string]any)
			if data, err := os.ReadFile(output); err == nil {
				if err := json.Unmarshal(data, &conf); err != nil {
					return fmt.Errorf("failed to parse existing %s: %w", output, err)
				}
			} else if !os.IsNotExist(err) {
				return fmt.Errorf("failed to read %s: %w", output, err)
			}

			// Ensure mcpServers key exists and is a map
			mcpServers, ok := conf["mcpServers"].(map[string]any)
			if !ok {
				mcpServers = make(map[string]any)
				conf["mcpServers"] = mcpServers
			}

			// Add/update geekbudget server config
			mcpServers["geekbudget"] = map[string]any{
				"command": execPath,
				"args":    []string{"mcp", "--username", username},
				"env": map[string]string{
					"GB_DBPATH": dbPath,
				},
			}

			// Write config file
			data, err := json.MarshalIndent(conf, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal config: %w", err)
			}

			if err := os.WriteFile(output, data, 0o644); err != nil {
				return fmt.Errorf("failed to write config file: %w", err)
			}

			absOutput, _ := filepath.Abs(output)
			fmt.Fprintf(os.Stderr, "MCP configuration written to: %s\n", absOutput)
			fmt.Fprintf(os.Stderr, "\nTo use with Claude Desktop or Claude Code:\n")
			fmt.Fprintf(os.Stderr, "1. Copy this file to your Claude configuration directory\n")
			fmt.Fprintf(os.Stderr, "2. Restart Claude\n")
			fmt.Fprintf(os.Stderr, "3. GeekBudget tools will be available for queries\n")

			return nil
		},
		Args: cobra.NoArgs,
	}

	res.Flags().StringVarP(&username, "username", "u", "", "username for data access")
	res.Flags().StringVarP(&output, "output", "o", ".mcp.json", "output file path")
	_ = res.MarkFlagRequired("username")

	return res
}
