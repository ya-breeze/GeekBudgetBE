package server

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

type AggregationsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewAggregationsAPIServiceImpl(logger *slog.Logger, db database.Storage,
) goserver.AggregationsAPIServicer {
	return &AggregationsAPIServiceImpl{logger: logger, db: db}
}

func (s *AggregationsAPIServiceImpl) GetBalances(context.Context, time.Time, time.Time, string,
) (goserver.ImplResponse, error) {
	return goserver.Response(500, nil), nil
}

func (s *AggregationsAPIServiceImpl) GetExpenses(
	ctx context.Context, dateFrom, dateTo time.Time, outputCurrencyID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts")
		return goserver.Response(500, nil), nil
	}

	transactions, err := s.db.GetTransactions(userID, dateFrom, dateTo)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions")
		return goserver.Response(500, nil), nil
	}

	res := Aggregate(accounts, transactions, dateFrom, dateTo, utils.GranularityMonth, s.logger)

	return goserver.Response(200, res), nil
}

//nolint:funlen,cyclop // TODO: refactor
func Aggregate(
	accounts []goserver.Account, transactions []goserver.Transaction,
	dateFrom, dateTo time.Time, granularity utils.Granularity,
	log *slog.Logger,
) goserver.Aggregation {
	if dateFrom.IsZero() {
		dateFrom = utils.RoundToGranularity(time.Now(), granularity, false)
	}
	if dateTo.IsZero() {
		dateTo = utils.RoundToGranularity(time.Now(), granularity, true)
	}

	res := goserver.Aggregation{
		From: utils.RoundToGranularity(dateFrom, granularity, false),
		To:   utils.RoundToGranularity(dateTo, granularity, true),
	}
	res.Intervals = getIntervals(res.From, res.To, granularity)

	res.Currencies = []goserver.CurrencyAggregation{}
	for _, t := range transactions {
		if t.Date.Before(res.From) || t.Date.After(res.To) {
			continue
		}
		intervalIdx := -1
		for i, interval := range res.Intervals {
			if t.Date.Before(interval) {
				intervalIdx = i - 1
			}
		}
		if intervalIdx < 0 {
			intervalIdx = len(res.Intervals) - 1
		}

		movements := getExpenseMovements(accounts, t)
		for _, m := range movements {
			currencyIdx := slices.IndexFunc(res.Currencies,
				func(item goserver.CurrencyAggregation) bool {
					return item.CurrencyId == m.CurrencyId
				})
			if currencyIdx == -1 {
				res.Currencies = append(res.Currencies, goserver.CurrencyAggregation{CurrencyId: m.CurrencyId})
				currencyIdx = len(res.Currencies) - 1
			}

			accountIdx := slices.IndexFunc(res.Currencies[currencyIdx].Accounts,
				func(item goserver.AccountAggregation) bool {
					return item.AccountId == m.AccountId
				})
			if accountIdx == -1 {
				res.Currencies[currencyIdx].Accounts = append(res.Currencies[currencyIdx].Accounts,
					goserver.AccountAggregation{
						AccountId: m.AccountId,
						Amounts:   make([]float64, len(res.Intervals)),
					})
				accountIdx = len(res.Currencies[currencyIdx].Accounts) - 1
			}

			res.Currencies[currencyIdx].Accounts[accountIdx].Amounts[intervalIdx] += m.Amount
		}
	}

	return res
}

func getExpenseMovements(accounts []goserver.Account, t goserver.Transaction) []goserver.Movement {
	movements := []goserver.Movement{}
	for _, m := range t.Movements {
		if isExpenseAccount(accounts, m.AccountId) {
			movements = append(movements, m)
		}
	}

	return movements
}

func isExpenseAccount(accounts []goserver.Account, accountID string) bool {
	for _, a := range accounts {
		if a.Id == accountID {
			if a.Type == constants.AccountExpense {
				return true
			}
		}
	}

	return false
}

func getIntervals(dateFrom, dateTo time.Time, granularity utils.Granularity,
) []time.Time {
	intervals := []time.Time{}
	for dateFrom.Before(dateTo) {
		intervals = append(intervals, dateFrom)
		switch granularity {
		case utils.GranularityMonth:
			dateFrom = dateFrom.AddDate(0, 1, 0)
		case utils.GranularityYear:
			dateFrom = dateFrom.AddDate(1, 0, 0)
		default:
			break
		}
	}
	return intervals
}

func (s *AggregationsAPIServiceImpl) GetIncomes(context.Context, time.Time, time.Time, string,
) (goserver.ImplResponse, error) {
	return goserver.Response(500, nil), nil
}
