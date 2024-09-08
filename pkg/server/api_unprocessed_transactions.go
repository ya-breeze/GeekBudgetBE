package server

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
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
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	matchers, err := s.db.GetMatchersRuntime(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matchers")
		return goserver.Response(500, nil), nil
	}

	transactions, err := s.db.GetTransactions(userID, time.Time{}, time.Time{})
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions")
		return goserver.Response(500, nil), nil
	}
	transactions = s.filterUnprocessedTransactions(transactions)

	res := make([]goserver.UnprocessedTransaction, 0, len(transactions))
	for _, t := range transactions {
		m, err := s.matchUnprocessedTransactions(matchers, t)
		if err != nil {
			s.logger.With("error", err).Error("Failed to match unprocessed transaction")
			return goserver.Response(500, nil), nil
		}

		res = append(res, goserver.UnprocessedTransaction{
			Transaction: t,
			Matched:     m,
			Duplicates:  nil,
		})
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Transaction.Date.Before(res[j].Transaction.Date)
	})
	return goserver.Response(200, res), nil
}

func (s *UnprocessedTransactionsAPIServiceImpl) ConvertUnprocessedTransaction(
	ctx context.Context,
	id string,
	transactionNoID goserver.TransactionNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}
	s.logger.Info("Converting unprocessed transaction", "transaction", id, "user", userID)

	transaction, err := s.db.UpdateTransaction(userID, id, &transactionNoID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to convert unprocessed transaction")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, transaction), nil
}

func (s *UnprocessedTransactionsAPIServiceImpl) DeleteUnprocessedTransaction(
	ctx context.Context,
	transactionID string,
	duplicateTransactionID string,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *UnprocessedTransactionsAPIServiceImpl) filterUnprocessedTransactions(transactions []goserver.Transaction,
) []goserver.Transaction {
	res := make([]goserver.Transaction, 0, len(transactions))
	for _, t := range transactions {
		for _, m := range t.Movements {
			if m.AccountId == "" {
				res = append(res, t)
				break
			}
		}
	}
	return res
}

func (s *UnprocessedTransactionsAPIServiceImpl) matchUnprocessedTransactions(
	matchers []database.MatcherRuntime, transactionSrc goserver.Transaction,
) ([]goserver.MatcherAndTransaction, error) {
	var transaction goserver.Transaction
	if err := utils.DeepCopy(&transactionSrc, &transaction); err != nil {
		return nil, fmt.Errorf("can't copy transaction: %w", err)
	}

	s.logger.With("transaction", transaction.Id).Info("Matching transaction")
	res := make([]goserver.MatcherAndTransaction, 0)

	for _, matcher := range matchers {
		if matcher.DescriptionRegexp != nil && !matcher.DescriptionRegexp.MatchString(transaction.Description) {
			continue
		}

		outputTransaction := models.TransactionWithoutID(&transaction)
		outputTransaction.Description = matcher.Matcher.OutputDescription
		for i := range outputTransaction.Movements {
			if outputTransaction.Movements[i].AccountId == "" {
				outputTransaction.Movements[i].AccountId = matcher.Matcher.OutputAccountId
			}
		}

		s.logger.With("matcher", matcher.Matcher.Id).Info("Matched transaction")
		res = append(res, goserver.MatcherAndTransaction{
			MatcherId:   matcher.Matcher.Id,
			Transaction: *outputTransaction,
		})
	}

	return res, nil
}
