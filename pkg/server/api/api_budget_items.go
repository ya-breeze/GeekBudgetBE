package api

import (
	"context"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
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
	// TODO: implement in next sub-task
	return goserver.Response(501, nil), nil
}

// CreateBudgetItem - create new budgetItem
func (s *BudgetItemsAPIServiceImpl) CreateBudgetItem(
	ctx context.Context, budgetItemNoID goserver.BudgetItemNoId,
) (goserver.ImplResponse, error) {
	// TODO: implement in next sub-task
	return goserver.Response(501, nil), nil
}

// GetBudgetItem - get budgetItem
func (s *BudgetItemsAPIServiceImpl) GetBudgetItem(ctx context.Context, id string) (goserver.ImplResponse, error) {
	// TODO: implement in next sub-task
	return goserver.Response(501, nil), nil
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
