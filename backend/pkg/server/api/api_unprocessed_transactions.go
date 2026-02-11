package api

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/shopspring/decimal"
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

func (s *UnprocessedTransactionsAPIServiceImpl) ProcessUnprocessedTransactionsAgainstMatcher(
	ctx context.Context, userID string, matcherID string, excludeTransactionID string,
) ([]string, error) {
	matcher, err := s.db.GetMatcher(userID, matcherID)
	if err != nil {
		s.logger.With("error", err, "matcherId", matcherID).Error("Failed to get matcher for auto-processing")
		return nil, err
	}

	// We only auto-process if the matcher is "perfect"
	// Perfect match defined as: at least 10 confirmations and all of them are true
	if len(matcher.ConfirmationHistory) < 10 {
		return nil, nil
	}
	for _, confirmed := range matcher.ConfirmationHistory {
		if !confirmed {
			return nil, nil
		}
	}

	// Get all matchers runtime to reuse the helper
	matchersRuntime, err := s.db.GetMatchersRuntime(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matchers runtime")
		return nil, err
	}

	// Filter down to just the specific matcher runtime we care about
	var specificMatcherRuntime *database.MatcherRuntime
	for i := range matchersRuntime {
		if matchersRuntime[i].Matcher.Id == matcherID {
			specificMatcherRuntime = &matchersRuntime[i]
			break
		}
	}
	if specificMatcherRuntime == nil {
		err := fmt.Errorf("matcher runtime not found for id %s", matcherID)
		s.logger.With("error", err).Error("Matcher runtime missing")
		return nil, err
	}

	// Get all transactions to find unprocessed ones and for duplicate checking
	// Optimization: This might be heavy if user has many transactions.
	// We might want to filter by date or similar in future, but for now we follow the pattern in PrepareUnprocessedTransactions
	allTransactions, err := s.db.GetTransactions(userID, time.Time{}, time.Time{}, false)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions")
		return nil, err
	}

	// Load accounts for filtering
	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts for auto-processing filtering")
		return nil, err
	}
	ignoreBeforeMap := make(map[string]time.Time)
	for _, acc := range accounts {
		if !acc.IgnoreUnprocessedBefore.IsZero() {
			ignoreBeforeMap[acc.Id] = acc.IgnoreUnprocessedBefore
		}
	}

	unprocessed := s.filterUnprocessedTransactions(allTransactions, ignoreBeforeMap)
	var processedIDs []string

	for _, t := range unprocessed {
		if t.Id == excludeTransactionID {
			continue
		}

		// Use the common matching logic
		matchDetails := common.MatchWithDetails(specificMatcherRuntime, &t)
		if !matchDetails.Matched {
			continue
		}

		// Conflict check: if it matches multiple matchers, don't auto-process
		matchesCount := 0
		for i := range matchersRuntime {
			if common.MatchWithDetails(&matchersRuntime[i], &t).Matched {
				matchesCount++
			}
		}
		if matchesCount > 1 {
			s.logger.With("transactionId", t.Id).Info("Skipping auto-processing due to multiple matcher matches")
			continue
		}

		// It matches! Let's convert it.
		// Construct the update payload
		transactionNoId := models.TransactionWithoutID(&t)
		transactionNoId.MatcherId = matcherID

		// 1. Prepare Proposed Movements to check if it would be a transfer
		proposedMovements := make([]goserver.Movement, len(transactionNoId.Movements))
		copy(proposedMovements, transactionNoId.Movements)
		for i := range proposedMovements {
			if proposedMovements[i].AccountId == "" {
				proposedMovements[i].AccountId = matcher.OutputAccountId
			}
		}

		// 2. Check for potential duplicate before auto-matching
		var duplicateFound *goserver.Transaction
		for _, existingT := range allTransactions {
			if existingT.Id == t.Id {
				continue
			}
			if common.IsDuplicate(transactionNoId.Date, transactionNoId.Movements, existingT.Date, existingT.Movements) {
				duplicateFound = &existingT
				break
			}
		}

		if duplicateFound != nil {
			transactionNoId.AutoMatchSkipReason = fmt.Sprintf("Potential duplicate detected: similar transaction exists from %s", duplicateFound.Date.Format("2006-01-02"))
			transactionNoId.IsAuto = false
			s.logger.With("transaction", t.Id, "duplicateId", duplicateFound.Id).Info("Skipping auto-processing due to duplicate")
		} else {
			// Apply matcher outputs
			description := matcher.OutputDescription
			tags := matcher.OutputTags
			if matcher.Simplified && matchDetails.MatchedKeyword != "" {
				description = matchDetails.MatchedOutput
				tags = append(append([]string{}, tags...), matchDetails.MatchedKeyword)
			}
			transactionNoId.Description = description
			transactionNoId.Tags = tags
			transactionNoId.Movements = proposedMovements
			transactionNoId.IsAuto = true
			// Merge tags
			transactionNoId.Tags = append(transactionNoId.Tags, matcher.OutputTags...)
			transactionNoId.Tags = sortAndRemoveDuplicates(transactionNoId.Tags)
		}

		// Persist
		_, err = s.db.UpdateTransaction(userID, t.Id, transactionNoId)
		if err != nil {
			s.logger.With("error", err, "transactionId", t.Id).Error("Failed to auto-process transaction")
			// We continue processing other transactions even if one fails
			continue
		}

		processedIDs = append(processedIDs, t.Id)
	}

	if len(processedIDs) > 0 {
		s.logger.Info("Auto-processed transactions", "count", len(processedIDs), "matcherId", matcherID)

		// Check balance after auto-processing
		// We need to know which accounts were affected. For simplicity, we check the matcher's account.
		if matcher.OutputAccountId != "" {
			if err := s.CheckBalanceForAccount(ctx, userID, matcher.OutputAccountId); err != nil {
				s.logger.With("error", err, "accountId", matcher.OutputAccountId).Error("Failed to check balance after auto-processing")
			}
		}
	}

	return processedIDs, nil
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
	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts for unprocessed filtering")
		return nil, 0, err
	}

	ignoreBeforeMap := make(map[string]time.Time)
	for _, acc := range accounts {
		if !acc.IgnoreUnprocessedBefore.IsZero() {
			ignoreBeforeMap[acc.Id] = acc.IgnoreUnprocessedBefore
		}
	}

	matchers, err := s.db.GetMatchersRuntime(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matchers")
		return nil, 0, err
	}

	var transactions []goserver.Transaction
	allTransactions, err := s.db.GetTransactions(userID, time.Time{}, time.Time{}, false)
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
	transactions = s.filterUnprocessedTransactions(transactions, ignoreBeforeMap)

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

	// compute all increases for the specified transaction (per currency)
	inc1 := common.GetIncreases(transaction.Movements)

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

		// skip transactions which are also unprocessed (have undefined accounts)
		isCandidateUnprocessed := false
		for _, m := range t.Movements {
			if m.AccountId == "" {
				isCandidateUnprocessed = true
				break
			}
		}
		if isCandidateUnprocessed {
			continue
		}

		// compute all increases in the transaction to compare (per currency)
		inc2 := common.GetIncreases(t.Movements)

		if len(inc1) != len(inc2) {
			continue
		}

		match := true
		for c, v1 := range inc1 {
			v2, ok := inc2[c]
			if !ok || v1.Sub(v2).Abs().GreaterThan(decimal.NewFromInt(1)) {
				match = false
				break
			}
		}

		if match {
			res = append(res, t)
		}
	}
	if len(res) != 0 {
		s.logger.Info("Found duplicates", "transaction", transaction.Id, "duplicates", len(res))
	}

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

func (s *UnprocessedTransactionsAPIServiceImpl) GetUnprocessedTransaction(
	ctx context.Context, id string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	transaction, err := s.db.GetTransaction(userID, id)
	if err != nil {
		if err == database.ErrNotFound {
			return goserver.Response(404, nil), nil
		}
		s.logger.With("error", err).Error("Failed to get transaction")
		return goserver.Response(500, nil), nil
	}

	matchers, err := s.db.GetMatchersRuntime(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matchers")
		return goserver.Response(500, nil), nil
	}

	m, err := s.matchUnprocessedTransactions(matchers, transaction)
	if err != nil {
		s.logger.With("error", err).Error("Failed to match unprocessed transaction")
		return goserver.Response(500, nil), nil
	}

	// Optimize duplicate search by time window +/- 2 days
	dateFrom := transaction.Date.Add(-48 * time.Hour)
	dateTo := transaction.Date.Add(48 * time.Hour)
	candidateTransactions, err := s.db.GetTransactions(userID, dateFrom, dateTo, false)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions for duplicate check")
		return goserver.Response(500, nil), nil
	}

	duplicates := s.getDuplicateTransactions(candidateTransactions, transaction)

	res := goserver.UnprocessedTransaction{
		Transaction: transaction,
		Matched:     m,
		Duplicates:  duplicates,
	}

	return goserver.Response(200, res), nil
}

func (s *UnprocessedTransactionsAPIServiceImpl) ConvertUnprocessedTransaction(
	ctx context.Context,
	id string,
	transactionNoID goserver.TransactionNoId,
	matcherId string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	var autoProcessedIds []string
	// If a matcher ID is provided, record a confirmation and update transaction
	if matcherId != "" {
		if err := s.db.AddMatcherConfirmation(userID, matcherId, true); err != nil {
			s.logger.With("error", err, "matcherId", matcherId).Error("Failed to add matcher confirmation")
			// We continue even if confirmation stats fail, as the conversion is the primary action
		} else {
			// Confirmation added successfully, now check if we can auto-process others
			ids, err := s.ProcessUnprocessedTransactionsAgainstMatcher(ctx, userID, matcherId, id)
			if err != nil {
				// Log but don't fail the request
				s.logger.With("error", err).Error("Failed to process unprocessed transactions against matcher")
			} else {
				autoProcessedIds = ids
			}
		}
		transactionNoID.MatcherId = matcherId
	}
	transactionNoID.IsAuto = false

	transaction, err := s.Convert(ctx, userID, id, &transactionNoID)
	if err != nil {
		return goserver.Response(500, nil), nil
	}

	// Check balance after conversion
	// We need to check all accounts moved in this transaction
	affectedAccounts := make(map[string]bool)
	for _, m := range transaction.Movements {
		if m.AccountId != "" {
			affectedAccounts[m.AccountId] = true
		}
	}
	for accID := range affectedAccounts {
		if err := s.CheckBalanceForAccount(ctx, userID, accID); err != nil {
			s.logger.With("error", err, "accountId", accID).Error("Failed to check balance after conversion")
		}
	}

	return goserver.Response(200, goserver.ConvertUnprocessedTransaction200Response{
		Transaction:      *transaction,
		AutoProcessedIds: autoProcessedIds,
	}), nil
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

func (s *UnprocessedTransactionsAPIServiceImpl) filterUnprocessedTransactions(
	transactions []goserver.Transaction, ignoreBeforeMap map[string]time.Time,
) []goserver.Transaction {
	res := make([]goserver.Transaction, 0, len(transactions))
	for _, t := range transactions {
		hasEmptyAccount := false
		shouldIgnore := false

		for _, m := range t.Movements {
			if m.AccountId == "" {
				hasEmptyAccount = true
			} else if ignoreDate, ok := ignoreBeforeMap[m.AccountId]; ok {
				// If any movement has a non-empty account which has ignore date set
				// and transaction is older than that date, we skip it.
				if t.Date.Before(ignoreDate) {
					shouldIgnore = true
					break
				}
			}
		}

		if hasEmptyAccount && !shouldIgnore {
			res = append(res, t)
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
		matchDetails := common.MatchWithDetails(&matcher, &transaction)
		if !matchDetails.Matched {
			continue
		}

		outputTransaction := models.TransactionWithoutID(&transaction)
		description := matcher.Matcher.OutputDescription
		tags := matcher.Matcher.OutputTags
		if matcher.Matcher.Simplified && matchDetails.MatchedKeyword != "" {
			description = matchDetails.MatchedOutput
			tags = append(append([]string{}, tags...), matchDetails.MatchedKeyword)
		}
		outputTransaction.Description = description
		outputTransaction.Tags = tags
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

func (s *UnprocessedTransactionsAPIServiceImpl) CheckBalanceForAccount(ctx context.Context, userID, accountID string) error {
	return common.CheckBalanceForAccount(ctx, s.logger, s.db, userID, accountID)
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
