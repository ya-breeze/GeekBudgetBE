package api

import (
	"context"
	"log/slog"
	"slices"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
	"github.com/ya-breeze/geekbudgetbe/pkg/utils"
)

type AggregationsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewAggregationsAPIServiceImpl(logger *slog.Logger, db database.Storage,
) *AggregationsAPIServiceImpl {
	return &AggregationsAPIServiceImpl{logger: logger, db: db}
}

func (s *AggregationsAPIServiceImpl) GetBalances(context.Context, time.Time, time.Time, string,
) (goserver.ImplResponse, error) {
	return goserver.Response(500, nil), nil
}

func (s *AggregationsAPIServiceImpl) GetExpenses(
	ctx context.Context, dateFrom, dateTo time.Time, outputCurrencyID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	aggregation, err := s.GetAggregatedExpenses(ctx, userID, dateFrom, dateTo, outputCurrencyID)
	if err != nil {
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, aggregation), nil
}

func (s *AggregationsAPIServiceImpl) GetAggregatedExpenses(
	ctx context.Context, userID string, dateFrom, dateTo time.Time, outputCurrencyID string,
) (*goserver.Aggregation, error) {
	if dateFrom.IsZero() {
		dateFrom = utils.RoundToGranularity(time.Now(), utils.GranularityMonth, false)
	}
	if dateTo.IsZero() {
		dateTo = utils.RoundToGranularity(time.Now(), utils.GranularityMonth, true)
	}

	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts")
		return nil, nil
	}

	transactions, err := s.db.GetTransactions(userID, dateFrom, dateTo)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions")
		return nil, nil
	}

	currenciesRatesFetcher := common.NewCurrenciesRatesFetcher(s.logger, s.db)
	res := Aggregate(
		accounts, transactions,
		dateFrom, dateTo,
		utils.GranularityMonth,
		outputCurrencyID, currenciesRatesFetcher,
		s.logger)

	return &res, nil
}

func Aggregate(
	accounts []goserver.Account, transactions []goserver.Transaction,
	dateFrom, dateTo time.Time, granularity utils.Granularity,
	outputCurrencyID string, currenciesRatesFetcher *common.CurrenciesRatesFetcher,
	log *slog.Logger,
) goserver.Aggregation {
	res := goserver.Aggregation{
		From: dateFrom,
		To:   dateTo,
	}
	res.Intervals = getIntervals(res.From, res.To, granularity)

	res.Currencies = []goserver.CurrencyAggregation{}
	for _, t := range transactions {
		if t.Date.Before(res.From) || t.Date.After(res.To) {
			log.Info("Ignore transaction", "date", t.Date)
			continue
		}
		intervalIdx := -1
		for i, interval := range res.Intervals {
			if t.Date.Before(interval) {
				intervalIdx = i - 1
				break
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
		if m.AccountId == "" || isExpenseAccount(accounts, m.AccountId) {
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
