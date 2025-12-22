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

func NewAggregationsAPIServiceImpl(logger *slog.Logger, db database.Storage) *AggregationsAPIServiceImpl {
	return &AggregationsAPIServiceImpl{
		logger: logger,
		db:     db,
	}
}

func (s *AggregationsAPIServiceImpl) GetIncomes(ctx context.Context, from time.Time, to time.Time, outputCurrencyID string, includeHidden bool) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	aggregation, err := s.GetAggregatedIncomes(ctx, userID, from, to, outputCurrencyID, includeHidden)
	if err != nil {
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, aggregation), nil
}

func (s *AggregationsAPIServiceImpl) GetAggregatedIncomes(
	ctx context.Context, userID string, dateFrom, dateTo time.Time, outputCurrencyID string, includeHidden bool,
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

	currencyMap := buildCurrencyMap(s.logger, s.db, userID)
	currenciesRatesFetcher := common.NewCurrenciesRatesFetcher(s.logger, s.db)

	res := Aggregate(
		ctx, accounts, transactions,
		dateFrom, dateTo,
		utils.GranularityMonth,
		outputCurrencyID, currenciesRatesFetcher,
		currencyMap,
		func(a goserver.Account) bool {
			return a.Type == constants.AccountIncome && (includeHidden || !a.HideFromReports)
		},
		s.logger)

	return &res, nil
}

func (s *AggregationsAPIServiceImpl) GetExpenses(
	ctx context.Context, dateFrom, dateTo time.Time, outputCurrencyID string, granularity string, includeHidden bool,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	aggGranularity := utils.GranularityMonth
	if granularity == "year" {
		aggGranularity = utils.GranularityYear
	}

	aggregation, err := s.GetAggregatedExpenses(ctx, userID, dateFrom, dateTo, outputCurrencyID, aggGranularity, includeHidden)
	if err != nil {
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, aggregation), nil
}

func (s *AggregationsAPIServiceImpl) GetAggregatedExpenses(
	ctx context.Context, userID string, dateFrom, dateTo time.Time, outputCurrencyID string, granularity utils.Granularity, includeHidden bool,
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
		granularity,
		outputCurrencyID, currenciesRatesFetcher,
		currencyMap,
		func(a goserver.Account) bool {
			return isExpenseAccount(a) && (includeHidden || !a.HideFromReports)
		},
		s.logger)

	return &res, nil
}

func (s *AggregationsAPIServiceImpl) GetBalances(
	ctx context.Context, dateFrom, dateTo time.Time, outputCurrencyID string, includeHidden bool,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	aggregation, err := s.GetAggregatedBalances(ctx, userID, dateFrom, dateTo, outputCurrencyID, includeHidden)
	if err != nil {
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, aggregation), nil
}

func (s *AggregationsAPIServiceImpl) GetAggregatedBalances(
	ctx context.Context, userID string, dateFrom, dateTo time.Time, outputCurrencyID string, includeHidden bool,
) (*goserver.Aggregation, error) {
	if dateFrom.IsZero() {
		// Use a very old date to capture all history
		dateFrom = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
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

	currencyMap := buildCurrencyMap(s.logger, s.db, userID)
	currenciesRatesFetcher := common.NewCurrenciesRatesFetcher(s.logger, s.db)

	filter := func(a goserver.Account) bool {
		return isAssetAccount(a) && (includeHidden || !a.HideFromReports)
	}

	res := Aggregate(
		ctx, accounts, transactions,
		dateFrom, dateTo,
		utils.GranularityMonth,
		outputCurrencyID, currenciesRatesFetcher,
		currencyMap,
		filter,
		s.logger)

	// Calculate initial balances (before dateFrom) to make the graph cumulative
	initialBalances, err := s.calculateInitialBalances(
		ctx, userID, accounts, dateFrom, outputCurrencyID, currencyMap, currenciesRatesFetcher, filter)
	if err != nil {
		s.logger.With("error", err).Error("Failed to calculate initial balances")
		return nil, nil
	}

	// Apply cumulative sum
	for i := range res.Currencies {
		currencyID := res.Currencies[i].CurrencyId
		for j := range res.Currencies[i].Accounts {
			accountID := res.Currencies[i].Accounts[j].AccountId

			// Get initial balance for this account and currency
			runningBalance := 0.0
			if accountBalances, ok := initialBalances[currencyID]; ok {
				runningBalance = accountBalances[accountID]
			}

			for k := range res.Currencies[i].Accounts[j].Amounts {
				runningBalance += res.Currencies[i].Accounts[j].Amounts[k]
				res.Currencies[i].Accounts[j].Amounts[k] = runningBalance
			}
		}
	}

	// There is a case where Account has initial balance but NO transactions in the selected period.
	// We need to add these accounts to the result as well, otherwise the line will be missing.
	for currencyID, accountBalances := range initialBalances {
		// Find or create currency in result
		currIdx := slices.IndexFunc(res.Currencies, func(c goserver.CurrencyAggregation) bool { return c.CurrencyId == currencyID })
		if currIdx == -1 {
			res.Currencies = append(res.Currencies, goserver.CurrencyAggregation{CurrencyId: currencyID, Accounts: []goserver.AccountAggregation{}})
			currIdx = len(res.Currencies) - 1
		}

		for accountID, initialAmount := range accountBalances {
			// Find or create account in result
			accIdx := slices.IndexFunc(res.Currencies[currIdx].Accounts, func(a goserver.AccountAggregation) bool { return a.AccountId == accountID })
			if accIdx == -1 {
				// If account wasn't in the result (no transactions in period), add it with flat line = initialAmount
				amounts := make([]float64, len(res.Intervals))
				for k := range amounts {
					amounts[k] = initialAmount
				}
				res.Currencies[currIdx].Accounts = append(res.Currencies[currIdx].Accounts, goserver.AccountAggregation{
					AccountId: accountID,
					Amounts:   amounts,
				})
			}
		}
	}

	return &res, nil
}

func (s *AggregationsAPIServiceImpl) calculateInitialBalances(
	ctx context.Context, userID string, accounts []goserver.Account, dateFrom time.Time,
	outputCurrencyID string, currencyMap map[string]string,
	currenciesRatesFetcher *common.CurrenciesRatesFetcher,
	filter AccountFilter,
) (map[string]map[string]float64, error) {
	// Map: CurrencyID -> AccountID -> Amount
	balances := make(map[string]map[string]float64)

	// 1. Sum up Opening Balances from Accounts
	for _, account := range accounts {
		if !filter(account) {
			continue
		}
		for _, balance := range account.BankInfo.Balances {
			// Convert opening balance to output currency
			movement := goserver.Movement{
				Amount:     balance.OpeningBalance,
				CurrencyId: balance.CurrencyId,
				AccountId:  account.Id,
			}
			// Reuse convertMovementAmount logic? It requires slightly different params but logic is same.
			// We use dateFrom as the reference date for conversion.
			convertedAmount, targetCurrencyID := convertMovementAmount(
				ctx, movement, dateFrom, outputCurrencyID, currencyMap[outputCurrencyID],
				currencyMap, currenciesRatesFetcher, s.logger)

			if _, ok := balances[targetCurrencyID]; !ok {
				balances[targetCurrencyID] = make(map[string]float64)
			}
			balances[targetCurrencyID][account.Id] += convertedAmount
		}
	}

	// 2. Sum up Past Transactions (if any)
	beginningOfTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	if dateFrom.After(beginningOfTime) {
		pastTransactions, err := s.db.GetTransactions(userID, beginningOfTime, dateFrom)
		if err != nil {
			return nil, err
		}

		// Aggregate past transactions into a single bucket
		// We use GranularityYear just to satisfy the function signature, but we'll sum up everything.
		resPast := Aggregate(
			ctx, accounts, pastTransactions,
			beginningOfTime, dateFrom,
			utils.GranularityYear,
			outputCurrencyID, currenciesRatesFetcher,
			currencyMap,
			filter,
			s.logger)

		// Merge past aggregation results into balances
		for _, currAgg := range resPast.Currencies {
			targetCurrencyID := currAgg.CurrencyId
			if _, ok := balances[targetCurrencyID]; !ok {
				balances[targetCurrencyID] = make(map[string]float64)
			}

			for _, accAgg := range currAgg.Accounts {
				total := 0.0
				for _, amount := range accAgg.Amounts {
					total += amount
				}
				balances[targetCurrencyID][accAgg.AccountId] += total
			}
		}
	}

	return balances, nil
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

type AccountFilter func(goserver.Account) bool

func Aggregate(
	ctx context.Context, accounts []goserver.Account, transactions []goserver.Transaction,
	dateFrom, dateTo time.Time, granularity utils.Granularity,
	outputCurrencyID string, currenciesRatesFetcher *common.CurrenciesRatesFetcher,
	currencyMap map[string]string,
	accountFilter AccountFilter,
	log *slog.Logger,
) goserver.Aggregation {
	// Ensure dateFrom starts at the beginning of a month/year to avoid interval drift
	// when adding months to dates like Jan 31 or Dec 31.
	dateFrom = utils.RoundToGranularity(dateFrom, granularity, false)
	// Round dateTo as well for consistency
	dateTo = utils.RoundToGranularity(dateTo, granularity, true)

	res := goserver.Aggregation{
		From: dateFrom,
		To:   dateTo,
	}
	res.Intervals = getIntervals(res.From, res.To, granularity)

	// Get the output currency name from the map if outputCurrencyID is provided
	outputCurrencyName := ""
	if outputCurrencyID != "" {
		outputCurrencyName = currencyMap[outputCurrencyID]
	}

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

		movements := getMovements(accounts, t, accountFilter)
		processMovements(ctx, movements, t.Date, outputCurrencyID, outputCurrencyName,
			currencyMap, currenciesRatesFetcher, intervalIdx, &res, log)
	}

	return res
}

func getMovements(accounts []goserver.Account, t goserver.Transaction, filter AccountFilter) []goserver.Movement {
	movements := []goserver.Movement{}
	for _, m := range t.Movements {
		if m.AccountId == "" || isAccountType(accounts, m.AccountId, filter) {
			movements = append(movements, m)
		}
	}

	return movements
}

func isAccountType(accounts []goserver.Account, accountID string, filter AccountFilter) bool {
	for _, a := range accounts {
		if a.Id == accountID {
			return filter(a)
		}
	}
	return false
}

func isExpenseAccount(a goserver.Account) bool {
	return a.Type == constants.AccountExpense
}

func isAssetAccount(a goserver.Account) bool {
	return a.Type == constants.AccountAsset
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
