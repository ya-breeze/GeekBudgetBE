package test

import (
	"log/slog"

	"github.com/dusted-go/logging/prettylog"
)

func CreateTestLogger() *slog.Logger {
	return slog.New(prettylog.NewHandler(&slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: false,
	}))
}
