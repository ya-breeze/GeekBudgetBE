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
	ctx context.Context, userID string, dateFrom, dateTo time.Time, outputCurrencyName string,
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

	// Prepare map currencyID->CurrencyName for all currencies of the current user.
	currencyMap := buildCurrencyMap(s.logger, s.db, userID)

	currenciesRatesFetcher := common.NewCurrenciesRatesFetcher(s.logger, s.db)
	res := Aggregate(
		ctx, accounts, transactions,
		dateFrom, dateTo,
		utils.GranularityMonth,
		outputCurrencyName, currenciesRatesFetcher,
		currencyMap,
		s.logger)

	return &res, nil
}

func buildCurrencyMap(logger *slog.Logger, storage database.Storage, userID string) map[string]string {
	currencies, err := storage.GetCurrencies(userID)
	if err != nil {
		logger.With("error", err, "userID", userID).Error("Failed to get currencies for user")
		return make(map[string]string)
	}

	currencyMap := make(map[string]string, len(currencies))
	for _, currency := range currencies {
		currencyMap[currency.Id] = currency.Name
	}

	logger.Debug("Built currency map", "userID", userID, "currencyCount", len(currencyMap))
	return currencyMap
}

func findCurrencyID(currencyName string, currencyMap map[string]string) string {
	if currencyName == "" {
		return ""
	}
	for id, name := range currencyMap {
		if name == currencyName {
			return id
		}
	}
	return ""
}

func processMovements(
	ctx context.Context, movements []goserver.Movement, transactionDate time.Time,
	outputCurrencyID, outputCurrencyName string, currencyMap map[string]string,
	currenciesRatesFetcher *common.CurrenciesRatesFetcher, intervalIdx int,
	res *goserver.Aggregation, log *slog.Logger,
) {
	for _, m := range movements {
		// Convert movement amount to target currency if needed
		convertedAmount, targetCurrencyID := convertMovementAmount(
			ctx, m, transactionDate, outputCurrencyID, outputCurrencyName,
			currencyMap, currenciesRatesFetcher, log)

		// Use target currency ID for grouping (either converted or original)
		currencyIdx := slices.IndexFunc(res.Currencies,
			func(item goserver.CurrencyAggregation) bool {
				return item.CurrencyId == targetCurrencyID
			})
		if currencyIdx == -1 {
			res.Currencies = append(res.Currencies, goserver.CurrencyAggregation{CurrencyId: targetCurrencyID})
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

		// Use converted amount for aggregation
		res.Currencies[currencyIdx].Accounts[accountIdx].Amounts[intervalIdx] += convertedAmount
	}
}

func Aggregate(
	ctx context.Context, accounts []goserver.Account, transactions []goserver.Transaction,
	dateFrom, dateTo time.Time, granularity utils.Granularity,
	outputCurrencyName string, currenciesRatesFetcher *common.CurrenciesRatesFetcher,
	currencyMap map[string]string,
	log *slog.Logger,
) goserver.Aggregation {
	res := goserver.Aggregation{
		From: dateFrom,
		To:   dateTo,
	}
	res.Intervals = getIntervals(res.From, res.To, granularity)

	// find ID of outputCurrencyMap
	outputCurrencyID := findCurrencyID(outputCurrencyName, currencyMap)

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
		processMovements(ctx, movements, t.Date, outputCurrencyID, outputCurrencyName,
			currencyMap, currenciesRatesFetcher, intervalIdx, &res, log)
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

func getIntervals(dateFrom, dateTo time.Time, granularity utils.Granularity) []time.Time {
	intervals := []time.Time{}
	for dateFrom.Before(dateTo) {
		intervals = append(intervals, dateFrom)
		switch granularity {
		case utils.GranularityMonth:
			dateFrom = dateFrom.AddDate(0, 1, 0)
		case utils.GranularityYear:
			dateFrom = dateFrom.AddDate(1, 0, 0)
		default:
			dateFrom = dateFrom.AddDate(1, 0, 0)
		}
	}
	return intervals
}

// convertMovementAmount converts a movement amount to the target currency using the provided exchange rate fetcher.
// Returns the converted amount and the target currency ID to use for grouping.
// If conversion is not needed or fails, returns the original amount and currency ID.
func convertMovementAmount(
	ctx context.Context, movement goserver.Movement, transactionDate time.Time,
	outputCurrencyID string, outputCurrencyName string,
	currencyMap map[string]string,
	currenciesRatesFetcher *common.CurrenciesRatesFetcher, log *slog.Logger,
) (float64, string) {
	// If no output currency specified, return original
	if outputCurrencyID == "" {
		log.Debug("No output currency specified, using original currency",
			"originalCurrency", movement.CurrencyId, "amount", movement.Amount)
		return movement.Amount, movement.CurrencyId
	}

	// If currencies rate fetcher is nil, return original
	if currenciesRatesFetcher == nil {
		log.Warn("Currency rate fetcher is nil, using original currency",
			"originalCurrency", movement.CurrencyId, "outputCurrency", outputCurrencyName, "amount", movement.Amount)
		return movement.Amount, movement.CurrencyId
	}

	// If same currency, no conversion needed
	if movement.CurrencyId == outputCurrencyID {
		log.Debug("Same currency, no conversion needed",
			"currency", movement.CurrencyId, "amount", movement.Amount)
		return movement.Amount, movement.CurrencyId
	}

	// Name of the original currency
	originalCurrencyName := currencyMap[movement.CurrencyId]

	// Attempt currency conversion
	convertedAmount, err := currenciesRatesFetcher.Convert(
		ctx, transactionDate, originalCurrencyName, outputCurrencyName, movement.Amount)
	if err != nil {
		log.Warn("Currency conversion failed, using original amount",
			"error", err,
			"date", transactionDate.Format("2006-01-02"),
			"fromCurrency", originalCurrencyName,
			"toCurrency", outputCurrencyName,
			"originalAmount", movement.Amount)
		return movement.Amount, movement.CurrencyId
	}

	log.Debug("Currency conversion successful",
		"date", transactionDate.Format("2006-01-02"),
		"fromCurrency", originalCurrencyName,
		"toCurrency", outputCurrencyName,
		"originalAmount", movement.Amount,
		"convertedAmount", convertedAmount)

	return convertedAmount, outputCurrencyID
}

func (s *AggregationsAPIServiceImpl) GetIncomes(context.Context, time.Time, time.Time, string,
) (goserver.ImplResponse, error) {
	return goserver.Response(500, nil), nil
}
