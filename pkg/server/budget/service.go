package budget

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
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

// ValidateFutureMonth ensures the given monthStart is in the future (not past or current month)
func (s *Service) ValidateFutureMonth(monthStart time.Time) error {
	now := time.Now()
	// Get first day of next month
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	if monthStart.Before(nextMonth) {
		return fmt.Errorf("budget month must be in the future, got %s but minimum is %s",
			monthStart.Format("2006-01"), nextMonth.Format("2006-01"))
	}

	return nil
}

// ListMonthlyBudget retrieves all budget items for a specific month and user
func (s *Service) ListMonthlyBudget(ctx context.Context, userID string, monthStart time.Time) ([]goserver.BudgetItem, error) {
	return s.db.GetBudgetItemsByMonth(userID, monthStart)
}
