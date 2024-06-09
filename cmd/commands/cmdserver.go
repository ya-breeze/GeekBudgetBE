package commands

import (
	"errors"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/server"
)

type ContextKey string

const ConfigKey ContextKey = "config"

func CmdServer() *cobra.Command {
	res := &cobra.Command{
		Use:   "server",
		Short: "Start HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			cfg, ok := cmd.Context().Value(ConfigKey).(*config.Config)
			if !ok {
				return errors.New("could not retrieve config from context")
			}
			return server.Server(logger, cfg)
		},
	}

	return res
}
