package api

import (
	"context"
	"errors"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

type MergedTransactionsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewMergedTransactionsAPIService(logger *slog.Logger, db database.Storage) goserver.MergedTransactionsAPIServicer {
	return &MergedTransactionsAPIServiceImpl{logger: logger, db: db}
}

func (s *MergedTransactionsAPIServiceImpl) GetMergedTransactions(ctx context.Context) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	merged, err := s.db.GetMergedTransactions(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get merged transactions")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, merged), nil
}

func (s *MergedTransactionsAPIServiceImpl) GetMergedTransaction(ctx context.Context, id string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	merged, err := s.db.GetMergedTransaction(userID, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return goserver.Response(404, nil), nil
		}
		s.logger.With("error", err, "id", id).Error("Failed to get merged transaction")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, merged), nil
}

func (s *MergedTransactionsAPIServiceImpl) UnmergeMergedTransaction(ctx context.Context, id string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	err := s.db.UnmergeTransaction(userID, id)
	if err != nil {
		s.logger.With("error", err, "id", id).Error("Failed to unmerge transaction")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, nil), nil
}
