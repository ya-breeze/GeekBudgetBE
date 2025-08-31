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

// SaveMonthlyBudget saves budget entries for a specific month, validating future date and expense accounts only
func (s *Service) SaveMonthlyBudget(ctx context.Context, userID string, monthStart time.Time, entries []goserver.BudgetItemNoId) error {
	// Validate future month
	if err := s.ValidateFutureMonth(monthStart); err != nil {
		return err
	}

	// Get all accounts for the user
	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts for budget validation")
		return fmt.Errorf("failed to get accounts: %w", err)
	}

	// Create map of expense account IDs for validation
	expenseAccountIDs := make(map[string]bool)
	for _, account := range accounts {
		if account.Type == "expense" {
			expenseAccountIDs[account.Id] = true
		}
	}

	// Validate that all entries are for expense accounts
	for _, entry := range entries {
		if !expenseAccountIDs[entry.AccountId] {
			return fmt.Errorf("budget entry for account %s is not an expense account", entry.AccountId)
		}
	}

	// Delete existing budget items for this month (for idempotency)
	existingItems, err := s.db.GetBudgetItemsByMonth(userID, monthStart)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get existing budget items")
		return fmt.Errorf("failed to get existing budget items: %w", err)
	}

	for _, item := range existingItems {
		if err := s.db.DeleteBudgetItem(userID, item.Id); err != nil {
			s.logger.With("error", err).Error("Failed to delete existing budget item")
			return fmt.Errorf("failed to delete existing budget item: %w", err)
		}
	}

	// Create new budget items
	for _, entry := range entries {
		// Set the date to the month start
		entry.Date = monthStart

		_, err := s.db.CreateBudgetItem(userID, &entry)
		if err != nil {
			s.logger.With("error", err).Error("Failed to create budget item")
			return fmt.Errorf("failed to create budget item: %w", err)
		}
	}

	s.logger.Info("Monthly budget saved", "userID", userID, "month", monthStart.Format("2006-01"), "entries", len(entries))
	return nil
}
