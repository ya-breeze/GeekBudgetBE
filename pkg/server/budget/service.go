package budget

import (
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
)

type Service struct {
	logger *slog.Logger
	db     database.Storage
}

func NewService(logger *slog.Logger, db database.Storage) *Service {
	return &Service{
		logger: logger,
		db:     db,
	}
}
