package api

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"sort"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

type UnprocessedTransactionsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewUnprocessedTransactionsAPIServiceImpl(logger *slog.Logger, db database.Storage,
) *UnprocessedTransactionsAPIServiceImpl {
	return &UnprocessedTransactionsAPIServiceImpl{logger: logger, db: db}
}

func (s *UnprocessedTransactionsAPIServiceImpl) Convert(
	ctx context.Context, userID string, id string, transactionNoID goserver.TransactionNoIdInterface,
) (*goserver.Transaction, error) {
	s.logger.Info("Converting unprocessed transaction", "transaction", id, "user", userID)

	transaction, err := s.db.UpdateTransaction(userID, id, transactionNoID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to convert unprocessed transaction")
		return nil, fmt.Errorf("failed to convert unprocessed transaction: %w", err)
	}

	return &transaction, nil
}

func (s *UnprocessedTransactionsAPIServiceImpl) PrepareUnprocessedTransactions(
	ctx context.Context, userID string, single bool, continuationID string,
) ([]goserver.UnprocessedTransaction, int, error) {
	matchers, err := s.db.GetMatchersRuntime(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matchers")
		return nil, 0, err
	}

	var transactions []goserver.Transaction
	allTransactions, err := s.db.GetTransactions(userID, time.Time{}, time.Time{})
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions")
		return nil, 0, err
	}
	transactions = allTransactions
	if len(continuationID) > 0 {
		for i, t := range transactions {
			if t.Id == continuationID {
				transactions = transactions[i+1:]
				break
			}
		}
	}
	transactions = s.filterUnprocessedTransactions(transactions)

	res := make([]goserver.UnprocessedTransaction, 0, len(transactions))
	for _, t := range transactions {
		m, err := s.matchUnprocessedTransactions(matchers, t)
		if err != nil {
			s.logger.With("error", err).Error("Failed to match unprocessed transaction")
			return nil, 0, err
		}

		d := s.getDuplicateTransactions(allTransactions, t)

		res = append(res, goserver.UnprocessedTransaction{
			Transaction: t,
			Matched:     m,
			Duplicates:  d,
		})

		if single {
			break
		}
	}

	return res, len(transactions), nil
}

func (s *UnprocessedTransactionsAPIServiceImpl) getDuplicateTransactions(
	transactions []goserver.Transaction, transaction goserver.Transaction,
) []goserver.Transaction {
	res := make([]goserver.Transaction, 0)

	// compute all increases for the specified transaction
	var increases float64
	for _, m := range transaction.Movements {
		if m.Amount > 0 {
			increases += m.Amount
		}
	}

	for _, t := range transactions {
		if t.Id == transaction.Id {
			continue
		}

		// skip transactions which didn't happen within 2 days
		delta := t.Date.Sub(transaction.Date)
		if delta < 0 {
			delta = -delta
		}
		if delta > 2*time.Hour*24 {
			continue
		}

		// compute all increases in the transaction to compare
		var d float64
		for _, m := range t.Movements {
			if m.Amount > 0 {
				d += m.Amount
			}
		}
		if math.Abs(increases-d) > 1 {
			continue
		}
		res = append(res, t)
	}
	s.logger.Info("Found duplicates", "transaction", transaction.Id, "duplicates", len(res))

	return res
}

func (s *UnprocessedTransactionsAPIServiceImpl) GetUnprocessedTransactions(
	ctx context.Context,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	res, _, err := s.PrepareUnprocessedTransactions(ctx, userID, false, "")
	if err != nil {
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, res), nil
}

func (s *UnprocessedTransactionsAPIServiceImpl) ConvertUnprocessedTransaction(
	ctx context.Context,
	id string,
	transactionNoID goserver.TransactionNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}
	transaction, err := s.Convert(ctx, userID, id, &transactionNoID)
	if err != nil {
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, transaction), nil
}

func (s *UnprocessedTransactionsAPIServiceImpl) Delete(
	ctx context.Context,
	userID string,
	transactionID string,
	duplicateTransactionID string,
) error {
	return s.db.DeleteDuplicateTransaction(userID, transactionID, duplicateTransactionID)
}

func (s *UnprocessedTransactionsAPIServiceImpl) DeleteUnprocessedTransaction(
	ctx context.Context,
	transactionID string,
	duplicateTransactionID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	err := s.Delete(ctx, userID, transactionID, duplicateTransactionID)
	if err != nil {
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(204, nil), nil
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

	res := make([]goserver.MatcherAndTransaction, 0)

	for _, matcher := range matchers {
		if matcher.DescriptionRegexp != nil && !matcher.DescriptionRegexp.MatchString(transaction.Description) {
			continue
		}

		if matcher.PartnerAccountRegexp != nil &&
			!matcher.PartnerAccountRegexp.MatchString(transaction.PartnerAccount) {
			continue
		}

		outputTransaction := models.TransactionWithoutID(&transaction)
		outputTransaction.Description = matcher.Matcher.OutputDescription
		for i := range outputTransaction.Movements {
			if outputTransaction.Movements[i].AccountId == "" {
				outputTransaction.Movements[i].AccountId = matcher.Matcher.OutputAccountId
			}
		}

		outputTransaction.Tags = append(outputTransaction.Tags, matcher.Matcher.OutputTags...)
		outputTransaction.Tags = sortAndRemoveDuplicates(outputTransaction.Tags)

		res = append(res, goserver.MatcherAndTransaction{
			MatcherId:   matcher.Matcher.Id,
			Transaction: *outputTransaction,
		})
	}

	return res, nil
}

func sortAndRemoveDuplicates(input []string) []string {
	uniqueMap := make(map[string]struct{})
	for _, str := range input {
		uniqueMap[str] = struct{}{}
	}

	uniqueList := make([]string, 0, len(uniqueMap))
	for key := range uniqueMap {
		uniqueList = append(uniqueList, key)
	}

	sort.Strings(uniqueList)
	return uniqueList
}
