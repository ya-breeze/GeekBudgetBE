package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/cmd/commands"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := &config.Config{}

	var rootCmd = &cobra.Command{
		Use:   "geekbudget",
		Short: "GeekBudget is a personal finance manager",
	}
	rootCmd.AddCommand(
		commands.CmdUser(),
		commands.CmdServer(logger, cfg),
	)
	if err := rootCmd.Execute(); err != nil {
		logger.Error("ERROR:", "error", err)
		os.Exit(1)
	}
	return
}
