package server

import (
	"context"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type UnprocessedTransactionsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewUnprocessedTransactionsAPIServiceImpl(logger *slog.Logger, db database.Storage,
) goserver.UnprocessedTransactionsAPIServicer {
	return &UnprocessedTransactionsAPIServiceImpl{logger: logger, db: db}
}

func (s *UnprocessedTransactionsAPIServiceImpl) GetUnprocessedTransactions(
	ctx context.Context,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *UnprocessedTransactionsAPIServiceImpl) ConvertUnprocessedTransaction(
	ctx context.Context,
	transactionID string,
	transaction goserver.TransactionNoId,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *UnprocessedTransactionsAPIServiceImpl) DeleteUnprocessedTransaction(
	ctx context.Context,
	transactionID string,
	duplicateTransactionID string,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}
