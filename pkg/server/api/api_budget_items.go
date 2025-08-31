package api

import (
	"context"
	"log/slog"

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
	// TODO: implement in next sub-task
	return goserver.Response(501, nil), nil
}

// DeleteBudgetItem - delete budgetItem
func (s *BudgetItemsAPIServiceImpl) DeleteBudgetItem(ctx context.Context, id string) (goserver.ImplResponse, error) {
	// TODO: implement in next sub-task
	return goserver.Response(501, nil), nil
}
