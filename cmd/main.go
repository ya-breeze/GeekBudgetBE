//nolint:forbidigo // it's okay to use fmt in this file
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/cmd/commands"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
)

func main() {
	var cfgFile string

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
		commands.CmdUser(),
		commands.CmdServer(),
		commands.CmdFio(),
	)
	if err := rootCmd.Execute(); err != nil {
		fmt.Printf("ERROR: %s", err)
		os.Exit(1)
	}
}
