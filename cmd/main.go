//nolint:forbidigo // it's okay to use fmt in this file
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/cmd/commands"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
)

// Write to stdout in gray color
type grayWriter struct {
}

func (w *grayWriter) Write(p []byte) (n int, err error) {
	return fmt.Fprint(os.Stdout, "\x1b[90m", string(p), "\x1b[0m")
}

func main() {
	var cfgFile string

	log := slog.New(slog.NewJSONHandler(&grayWriter{}, nil))

	rootCmd := &cobra.Command{
		Use:   "geekbudget",
		Short: "GeekBudget is a personal finance manager",
		PersistentPreRun: func(cmd *cobra.Command, _ []string) {
			cfg, err := config.InitiateConfig(cfgFile)
			if err != nil {
				fmt.Printf("ERROR: %s", err)
				os.Exit(1)
			}
			cmd.SetContext(context.WithValue(cmd.Context(), commands.ConfigKey, cfg))
		},
	}

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
	rootCmd.AddCommand(
		commands.CmdUser(log),
		commands.CmdServer(),
		commands.CmdFio(log),
		commands.CmdRevolut(log),
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
