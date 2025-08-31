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

// Row represents a single account's budget vs actual comparison
type Row struct {
	AccountID   string  `json:"accountId"`
	AccountName string  `json:"accountName"`
	Planned     float64 `json:"planned"`
	Actual      float64 `json:"actual"`
	Delta       float64 `json:"delta"` // Actual - Planned
}

// Comparison represents the complete budget vs actual comparison for a month
type Comparison struct {
	Rows         []Row   `json:"rows"`
	TotalPlanned float64 `json:"totalPlanned"`
	TotalActual  float64 `json:"totalActual"`
	TotalDelta   float64 `json:"totalDelta"`
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
func (s *Service) ListMonthlyBudget(
	ctx context.Context, userID string, monthStart time.Time,
) ([]goserver.BudgetItem, error) {
	return s.db.GetBudgetItemsByMonth(userID, monthStart)
}

// SaveMonthlyBudget saves budget entries for a specific month; validates expense accounts only
func (s *Service) SaveMonthlyBudget(ctx context.Context, userID string, monthStart time.Time, entries []goserver.BudgetItemNoId) error {
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

// CompareMonthly compares planned budget vs actual expenses for a specific month
func (s *Service) CompareMonthly(ctx context.Context, userID string, monthStart time.Time, outputCurrencyName string) (*Comparison, error) {
	// Get planned amounts from budget items
	budgetItems, err := s.db.GetBudgetItemsByMonth(userID, monthStart)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get budget items for comparison")
		return nil, fmt.Errorf("failed to get budget items: %w", err)
	}

	// Build planned map
	plannedMap := make(map[string]float64)
	for _, item := range budgetItems {
		plannedMap[item.AccountId] = item.Amount
	}

	// Get all accounts (expense only for budgeting)
	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts for comparison")
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	// Get actual expenses for the month
	monthEnd := monthStart.AddDate(0, 1, 0)
	transactions, err := s.db.GetTransactions(userID, monthStart, monthEnd)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions for comparison")
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	// Calculate actual expenses per account
	actualMap := make(map[string]float64)
	for _, transaction := range transactions {
		for _, movement := range transaction.Movements {
			// Only count expense movements (positive amounts for expense accounts)
			if movement.Amount > 0 {
				actualMap[movement.AccountId] += movement.Amount
			}
		}
	}

	// Build comparison rows
	var rows []Row
	var totalPlanned, totalActual float64

	for _, account := range accounts {
		if account.Type == "expense" {
			planned := plannedMap[account.Id]
			actual := actualMap[account.Id]
			delta := actual - planned

			if planned > 0 || actual > 0 { // Only include accounts with activity
				rows = append(rows, Row{
					AccountID:   account.Id,
					AccountName: account.Name,
					Planned:     planned,
					Actual:      actual,
					Delta:       delta,
				})
			}

			totalPlanned += planned
			totalActual += actual
		}
	}

	return &Comparison{
		Rows:         rows,
		TotalPlanned: totalPlanned,
		TotalActual:  totalActual,
		TotalDelta:   totalActual - totalPlanned,
	}, nil
}

// CopyFromPreviousMonth copies budget items from a previous month to a new month
func (s *Service) CopyFromPreviousMonth(
	ctx context.Context, userID string, fromMonthStart, toMonthStart time.Time,
) (int, error) {
	// Fetch target month existing items
	existingItems, err := s.db.GetBudgetItemsByMonth(userID, toMonthStart)
	if err != nil {
		s.logger.With("error", err).Error("Failed to check existing budget items")
		return 0, fmt.Errorf("failed to check existing budget items: %w", err)
	}

	// Build map of target amounts by account
	targetMap := make(map[string]float64, len(existingItems))
	for _, it := range existingItems {
		targetMap[it.AccountId] = it.Amount
	}

	// Fetch source items
	sourceItems, err := s.db.GetBudgetItemsByMonth(userID, fromMonthStart)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get source budget items")
		return 0, fmt.Errorf("failed to get source budget items: %w", err)
	}

	count := 0
	for _, src := range sourceItems {
		amt, exists := targetMap[src.AccountId]
		if !exists || amt == 0 {
			item := goserver.BudgetItemNoId{
				Date:        toMonthStart,
				AccountId:   src.AccountId,
				Amount:      src.Amount,
				Description: src.Description,
			}
			if _, err := s.db.CreateBudgetItem(userID, &item); err != nil {
				s.logger.With("error", err).Error("Failed to create budget item during copy")
				return 0, fmt.Errorf("failed to create budget item during copy: %w", err)
			}
			count++
		}
	}

	s.logger.Info("Budget items copied",
		"userID", userID,
		"from", fromMonthStart.Format("2006-01"),
		"to", toMonthStart.Format("2006-01"),
		"count", count)
	return count, nil
}
