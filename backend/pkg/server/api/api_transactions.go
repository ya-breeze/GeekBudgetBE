package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

type TransactionsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewTransactionsAPIService(logger *slog.Logger, db database.Storage) goserver.TransactionsAPIServicer {
	return &TransactionsAPIServiceImpl{logger: logger, db: db}
}

func (s *TransactionsAPIServiceImpl) GetTransactions(
	ctx context.Context,
	descriptionParam string,
	amountFromParam, amountToParam float64,
	dateFrom, dateTo time.Time,
	onlySuspicious bool,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	transactions, err := s.db.GetTransactions(userID, dateFrom, dateTo, onlySuspicious)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, transactions), nil
}

func (s *TransactionsAPIServiceImpl) CreateTransaction(
	ctx context.Context, transactionNoID goserver.TransactionNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}
	s.logger.Info("Processing transaction create", "user", userID)

	transaction, err := s.db.CreateTransaction(userID, &transactionNoID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create transaction")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, transaction), nil
}

func (s *TransactionsAPIServiceImpl) UpdateTransaction(
	ctx context.Context, transactionID string, transactionNoID goserver.TransactionNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}
	s.logger.Info("Processing transaction update", "transaction", transactionID, "user", userID)

	transaction, err := s.db.UpdateTransaction(userID, transactionID, &transactionNoID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update transaction")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, transaction), nil
}

func (s *TransactionsAPIServiceImpl) DeleteTransaction(
	ctx context.Context, transactionID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	err := s.db.DeleteTransaction(userID, transactionID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to delete transaction")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, nil), nil
}

func (s *TransactionsAPIServiceImpl) GetTransaction(
	ctx context.Context, transactionID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	transaction, err := s.db.GetTransaction(userID, transactionID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transaction")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, transaction), nil
}

func (s *TransactionsAPIServiceImpl) MergeTransactions(
	ctx context.Context, mergeRequest goserver.MergeTransactionsRequest,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	s.logger.Info("Processing transactions merge", "keep", mergeRequest.KeepId, "merge", mergeRequest.MergeId, "user", userID)

	transaction, err := s.db.MergeTransactions(userID, mergeRequest.KeepId, mergeRequest.MergeId)
	if err != nil {
		s.logger.With("error", err).Error("Failed to merge transactions")
		return goserver.Response(400, nil), nil
	}

	return goserver.Response(200, transaction), nil
}
