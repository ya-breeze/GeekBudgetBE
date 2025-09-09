package api

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

type BudgetItemsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewBudgetItemsAPIService(logger *slog.Logger, db database.Storage) goserver.BudgetItemsAPIServicer {
	return &BudgetItemsAPIServiceImpl{
		logger: logger,
		db:     db,
	}
}

// validateFutureDate ensures the budget item date is in the future (not past or current month)
func (s *BudgetItemsAPIServiceImpl) validateFutureDate(date time.Time) error {
	now := time.Now()
	nextMonth := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, now.Location())

	if date.Before(nextMonth) {
		return fmt.Errorf("budget date must be in the future, got %s but minimum is %s",
			date.Format("2006-01"), nextMonth.Format("2006-01"))
	}

	return nil
}

// GetBudgetItems - get all budgetItems
func (s *BudgetItemsAPIServiceImpl) GetBudgetItems(ctx context.Context) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	budgetItems, err := s.db.GetBudgetItems(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get budget items")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, budgetItems), nil
}

// CreateBudgetItem - create new budgetItem
func (s *BudgetItemsAPIServiceImpl) CreateBudgetItem(
	ctx context.Context, budgetItemNoID goserver.BudgetItemNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Validate future date
	if err := s.validateFutureDate(budgetItemNoID.Date); err != nil {
		s.logger.With("error", err).Error("Budget item date validation failed")
		return goserver.Response(400, map[string]string{"error": err.Error()}), nil
	}

	budgetItem, err := s.db.CreateBudgetItem(userID, &budgetItemNoID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create budget item")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, budgetItem), nil
}

// GetBudgetItem - get budgetItem
func (s *BudgetItemsAPIServiceImpl) GetBudgetItem(ctx context.Context, id string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	budgetItem, err := s.db.GetBudgetItem(userID, id)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get budget item")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, budgetItem), nil
}

// UpdateBudgetItem - update budgetItem
func (s *BudgetItemsAPIServiceImpl) UpdateBudgetItem(
	ctx context.Context, id string, budgetItemNoID goserver.BudgetItemNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Validate future date
	if err := s.validateFutureDate(budgetItemNoID.Date); err != nil {
		s.logger.With("error", err).Error("Budget item date validation failed")
		return goserver.Response(400, map[string]string{"error": err.Error()}), nil
	}

	budgetItem, err := s.db.UpdateBudgetItem(userID, id, &budgetItemNoID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update budget item")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, budgetItem), nil
}

// DeleteBudgetItem - delete budgetItem
func (s *BudgetItemsAPIServiceImpl) DeleteBudgetItem(ctx context.Context, id string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	err := s.db.DeleteBudgetItem(userID, id)
	if err != nil {
		s.logger.With("error", err).Error("Failed to delete budget item")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, nil), nil
}
