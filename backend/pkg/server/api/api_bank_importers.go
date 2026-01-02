package api

import (
	"context"
	"fmt"
	"log/slog"
	"math"
	"os"
	"slices"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/bankimporters"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

type BankImportersAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewBankImportersAPIServiceImpl(
	logger *slog.Logger, db database.Storage,
) *BankImportersAPIServiceImpl {
	return &BankImportersAPIServiceImpl{logger: logger, db: db}
}

func (s *BankImportersAPIServiceImpl) CreateBankImporter(ctx context.Context, input goserver.BankImporterNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	res, err := s.db.CreateBankImporter(userID, &input)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create BankImporter")
		return goserver.Response(500, nil), nil
	}

	if forcedImports := common.GetForcedImportChannel(ctx); forcedImports != nil {
		s.logger.Info("Triggering forced import for new bank importer", "userID", userID, "bankImporterID", res.Id)
		forcedImports <- common.ForcedImport{
			UserID:         userID,
			BankImporterID: res.Id,
		}
	}

	return goserver.Response(200, res), nil
}

func (s *BankImportersAPIServiceImpl) DeleteBankImporter(ctx context.Context, id string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	err := s.db.DeleteBankImporter(userID, id)
	if err != nil {
		s.logger.With("error", err).Error("Failed to delete BankImporter")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, nil), nil
}

func (s *BankImportersAPIServiceImpl) GetBankImporters(ctx context.Context) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	bankImporters, err := s.db.GetBankImporters(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get bank importers")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, bankImporters), nil
}

func (s *BankImportersAPIServiceImpl) UpdateBankImporter(
	ctx context.Context, id string, input goserver.BankImporterNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	res, err := s.db.UpdateBankImporter(userID, id, &input)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update BankImporter")
		return goserver.Response(500, nil), nil
	}

	if forcedImports := common.GetForcedImportChannel(ctx); forcedImports != nil {
		s.logger.Info("Triggering forced import for updated bank importer", "userID", userID, "bankImporterID", id)
		forcedImports <- common.ForcedImport{
			UserID:         userID,
			BankImporterID: id,
		}
	}

	return goserver.Response(200, res), nil
}

func (s *BankImportersAPIServiceImpl) Fetch(
	ctx context.Context, userID, importerID string, isInteractive bool,
) (*goserver.ImportResult, error) {
	s.logger.Info("Fetching bank importer", "userID", userID, "bankImporterID", importerID)
	info, transactions, wasFetchAll, err := s.fetchFioTransactions(ctx, userID, importerID, isInteractive)
	if err != nil {
		s.logger.With("error", err).Error("Failed to fetch for bank importer")
		// Log failed import
		_ = s.addImportResult(userID, importerID, goserver.ImportResult{
			Date:        time.Now(),
			Status:      "error",
			Description: err.Error(),
		})
		return nil, err
	}

	lastImport, err := s.saveImportedTransactions(userID, importerID, info, transactions, wasFetchAll)
	if err != nil {
		s.logger.With("error", err).Error("Failed to save imported transactions")
		return nil, err
	}
	s.logger.Info("Bank importer fetched", "userID", userID, "result", lastImport)

	return lastImport, nil
}

func (s *BankImportersAPIServiceImpl) FetchBankImporter(
	ctx context.Context, id string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	lastImport, err := s.Fetch(ctx, userID, id, true)
	if err != nil {
		s.logger.With("error", err).Error("Failed to fetch")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, lastImport), nil
}

func (s *BankImportersAPIServiceImpl) fetchFioTransactions(
	ctx context.Context, userID, id string, isInteractive bool,
) (*goserver.BankAccountInfo, []goserver.TransactionNoId, bool, error) {
	s.logger.With("user", userID).Info("Fetching transactions for bank importer")

	biData, err := s.db.GetBankImporter(userID, id)
	if err != nil {
		return nil, nil, false, fmt.Errorf("can't fetch bank importer: %w", err)
	}

	if !isInteractive && biData.IsStopped {
		s.logger.With("userID", userID, "bankImporterID", id).Info("Bank importer is stopped, skipping fetch")
		return nil, nil, false, fmt.Errorf("bank importer is stopped")
	}

	currencies, err := s.db.GetCurrencies(userID)
	if err != nil {
		return nil, nil, false, fmt.Errorf("can't fetch currencies: %w", err)
	}

	cp := bankimporters.NewDefaultCurrencyProvider(s.db, userID, currencies)

	bi, err := bankimporters.NewFioConverter(s.logger, biData, cp)
	if err != nil {
		return nil, nil, false, fmt.Errorf("can't create FioConverter: %w", err)
	}

	s.logger.Info("Importing transactions")
	info, transactions, err := bi.Import(ctx)
	wasFetchAll := biData.FetchAll

	if err != nil {
		// Stop fetching if it was not interactive and fetch all is false
		if !biData.FetchAll && !isInteractive {
			s.logger.With("bankImporterID", id).With("userID", userID).Info("Bank importer failed, stopping further fetches")
			biData.IsStopped = true
			_, updateErr := s.db.UpdateBankImporter(userID, id, &biData)
			if updateErr != nil {
				s.logger.With("error", updateErr).Error("Failed to set IsStopped after import failure")
			}
		}

		if biData.FetchAll {
			s.logger.With("bankImporterID", id).With("userID", userID).Info("All transactions fetch failed. Disabling FetchAll and creating notification")
			biData.FetchAll = false
			_, updateErr := s.db.UpdateBankImporter(userID, id, &biData)
			if updateErr != nil {
				s.logger.With("error", updateErr).Error("Failed to reset FetchAll after import failure")
			}

			// Create notification
			_, notifyErr := s.db.CreateNotification(userID, &goserver.Notification{
				Date:        time.Now(),
				Type:        string(models.NotificationTypeError),
				Title:       "Bank Import Failed",
				Description: fmt.Sprintf("Failed to fetch all transactions for %q. The 'Fetch All' flag has been reset. Error: %s", biData.Name, err),
			})
			if notifyErr != nil {
				s.logger.With("error", notifyErr).Error("Failed to create notification for import failure")
			}
		} else if !isInteractive {
			// Create notification for stopped importer
			_, notifyErr := s.db.CreateNotification(userID, &goserver.Notification{
				Date:        time.Now(),
				Type:        string(models.NotificationTypeError),
				Title:       "Bank Import Stopped",
				Description: fmt.Sprintf("Failed to fetch transactions for %q. Automatic fetching has been stopped. Please check the importer settings. Error: %s", biData.Name, err),
			})
			if notifyErr != nil {
				s.logger.With("error", notifyErr).Error("Failed to create notification for stopped importer")
			}
		}

		return nil, nil, false, fmt.Errorf("can't import transactions: %w", err)
	}
	s.logger.With("info", info, "transactions", len(transactions)).Info("Imported transactions")

	if biData.IsStopped {
		s.logger.With("bankImporterID", id).With("userID", userID).Info("Bank importer fetched successfully, resetting IsStopped")
		biData.IsStopped = false
		_, err = s.db.UpdateBankImporter(userID, id, &biData)
		if err != nil {
			s.logger.With("error", err).Error("Failed to reset IsStopped after successful import")
		}
	}

	if biData.FetchAll {
		s.logger.With("bankImporterID", id).With("userID", userID).Info("All transactions fetched. Disabling FetchAll")
		biData.FetchAll = false
		_, err = s.db.UpdateBankImporter(userID, id, &biData)
		if err != nil {
			return nil, nil, false, fmt.Errorf("can't update BankImporter: %w", err)
		}
	}

	return info, transactions, wasFetchAll, nil
}

func (s *BankImportersAPIServiceImpl) updateLastImportFields(
	userID string, id string, info *goserver.BankAccountInfo, totalTransactionsCnt int, newTransactionsCnt int, suspiciousCnt int,
) (*goserver.ImportResult, error) {
	balances := []goserver.ImportResultBalancesInner{}
	if info != nil {
		for _, b := range info.Balances {
			balances = append(balances, goserver.ImportResultBalancesInner{
				Amount:     float32(b.ClosingBalance),
				CurrencyId: b.CurrencyId,
			})
		}
	}

	lastImport := goserver.ImportResult{
		Date:   time.Now(),
		Status: "success",
		Description: fmt.Sprintf("Fetched %d transactions. Imported %d new transactions.",
			totalTransactionsCnt, newTransactionsCnt),
		Balances:        balances,
		SuspiciousCount: int32(suspiciousCnt),
	}

	if err := s.addImportResult(userID, id, lastImport); err != nil {
		return nil, err
	}

	return &lastImport, nil
}

func (s *BankImportersAPIServiceImpl) addImportResult(
	userID, id string, result goserver.ImportResult,
) error {
	biData, err := s.db.GetBankImporter(userID, id)
	if err != nil {
		return fmt.Errorf("can't fetch bank importer: %w", err)
	}

	if result.Date.IsZero() {
		result.Date = time.Now()
	}

	if result.Status == "success" {
		biData.LastSuccessfulImport = result.Date
	}
	biData.LastImports = append(biData.LastImports, result)
	if len(biData.LastImports) > 10 {
		biData.LastImports = biData.LastImports[1:]
	}
	_, err = s.db.UpdateBankImporter(userID, id, &biData)
	if err != nil {
		return fmt.Errorf("can't update BankImporter: %w", err)
	}
	return nil
}

func (s *BankImportersAPIServiceImpl) saveImportedTransactions(
	userID, id string, info *goserver.BankAccountInfo, transactions []goserver.TransactionNoId, checkMissing bool,
) (*goserver.ImportResult, error) {
	if len(transactions) == 0 {
		return s.updateLastImportFields(userID, id, info, 0, 0, 0)
	}

	// Calculate the date range for fetching existing transactions
	// We want to fetch transactions starting from the earliest date in the import batch minus a margin
	earliestDate := transactions[0].Date
	for _, t := range transactions {
		if t.Date.Before(earliestDate) {
			earliestDate = t.Date
		}
	}
	var fetchFrom time.Time
	if !checkMissing {
		fetchFrom = earliestDate.AddDate(0, 0, -7) // 7 days safety margin
	}

	// Fetch all transactions from the database (including deleted ones)
	dbTransactions, err := s.db.GetTransactionsIncludingDeleted(userID, fetchFrom, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("can't fetch transactions from DB: %w", err)
	}

	matchers, err := s.db.GetMatchersRuntime(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matchers")
		return nil, fmt.Errorf("can't get matchers: %w", err)
	}

	biData, err := s.db.GetBankImporter(userID, id)
	if err != nil {
		return nil, fmt.Errorf("can't get bank importer: %w", err)
	}

	// Keep track of visited transactions within this batch to handle self-duplicates
	visitedExternalIDs := make(map[string]bool)

	// save transactions to the database
	cnt := 0
	for _, t := range transactions {

		// Imported transactions should have at least one external ID filled by the bank importer.
		// Revolut importer now initiates 2 IDs (legacy hash + stable hash)
		if len(t.ExternalIds) == 0 {
			return nil, fmt.Errorf("transaction has invalid external IDs: %v", t)
		}

		// search for existing transaction with the same external ID. If found, skip saving
		found := false

		// 0. Check against already visited in this batch (Self-Deduplication)
		for _, extID := range t.ExternalIds {
			if visitedExternalIDs[extID] {
				found = true
				s.logger.With("externalID", extID).Info("Duplicate transaction within import batch (external ID)")
				break
			}
		}

		if found {
			continue
		}

		// Check against DB transactions
		for _, dbt := range dbTransactions {
			// 1. Check for exact match of any external ID
			for _, extID := range t.ExternalIds {
				if slices.Contains(dbt.ExternalIds, extID) {
					found = true
					s.logger.With("externalID", extID).Info("Transaction already was imported (exact match)")
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			// Mark as visited so we don't process potential subsequent duplicates of this one in the same batch
			for _, extID := range t.ExternalIds {
				visitedExternalIDs[extID] = true
			}
			continue
		}

		// Try to match with perfect matchers
		// Create temporary transaction for matching
		tempDetails := &goserver.Transaction{
			Date:               t.Date,
			Description:        t.Description,
			Place:              t.Place,
			Tags:               t.Tags,
			PartnerName:        t.PartnerName,
			PartnerAccount:     t.PartnerAccount,
			PartnerInternalId:  t.PartnerInternalId,
			Extra:              t.Extra,
			UnprocessedSources: t.UnprocessedSources,
			ExternalIds:        t.ExternalIds,
			Movements:          t.Movements,
		}

		// Find all matchers that match
		var matches []struct {
			matcher      database.MatcherRuntime
			matchDetails common.MatchDetails
		}
		for _, matcher := range matchers {
			matchDetails := common.MatchWithDetails(&matcher, tempDetails)
			if matchDetails.Matched {
				matches = append(matches, struct {
					matcher      database.MatcherRuntime
					matchDetails common.MatchDetails
				}{matcher, matchDetails})
			}
		}

		// Only auto-process if exactly ONE matcher matches and it's a perfect match
		if len(matches) == 1 && isPerfectMatch(matches[0].matcher.Matcher) {
			matcher := matches[0].matcher
			matchDetails := matches[0].matchDetails

			s.logger.Info("Found perfect match", "matcher", matcher.Matcher.OutputDescription, "transaction", t.Description)

			// Apply matcher outputs
			description := matcher.Matcher.OutputDescription
			tags := matcher.Matcher.OutputTags
			if matcher.Matcher.Simplified && matchDetails.MatchedKeyword != "" {
				description = matchDetails.MatchedOutput
				tags = append(append([]string{}, tags...), matchDetails.MatchedKeyword)
			}
			t.Description = description
			for i := range t.Movements {
				if t.Movements[i].AccountId == "" {
					t.Movements[i].AccountId = matcher.Matcher.OutputAccountId
				}
			}
			t.Tags = append(t.Tags, matcher.Matcher.OutputTags...)
			t.Tags = sortAndRemoveDuplicates(t.Tags)
			t.MatcherId = matcher.Matcher.Id
			t.IsAuto = true

			// auto-confirm the matcher
			if err := s.db.AddMatcherConfirmation(userID, t.MatcherId, true); err != nil {
				s.logger.Warn("Failed to add confirmation to matcher", "matcher_id", t.MatcherId, "error", err)
			}
		}

		_, err = s.db.CreateTransaction(userID, &t)
		if err != nil {
			return nil, fmt.Errorf("can't save transaction: %w", err)
		}
		s.logger.Info("Imported transaction saved to DB")

		// Mark as visited
		for _, extID := range t.ExternalIds {
			visitedExternalIDs[extID] = true
		}

		cnt++
	}
	s.logger.With("count", cnt).Info("All new imported transactions saved to DB")

	suspiciousCnt := 0
	if checkMissing {
		// Create lookup map for incoming transactions (external IDs only)
		incomingExternalIDs := make(map[string]bool)

		for _, t := range transactions {
			for _, extID := range t.ExternalIds {
				incomingExternalIDs[extID] = true
			}
		}

		for _, dbt := range dbTransactions {
			// Check if transaction belongs to this account
			isAccountMatch := false
			for _, m := range dbt.Movements {
				if m.AccountId == biData.AccountId {
					isAccountMatch = true
					break
				}
			}
			if !isAccountMatch {
				continue
			}

			// Check if it exists in incoming transactions
			found := false
			// 1. External ID check
			for _, extID := range dbt.ExternalIds {
				if incomingExternalIDs[extID] {
					found = true
					break
				}
			}
			if found {
				continue
			}

			// 2. Duplicate check: if some transaction doesn't exist in fetched/uploaded list
			// then BE should also check list of duplicates of the existing transaction.
			// If any of them is present in fetch/upload list then transaction is not suspicious.
			if !found {
				for _, t := range transactions {
					if isDuplicate(&t, &dbt) {
						found = true
						s.logger.With("dbtID", dbt.Id, "incomingDesc", t.Description).Info("Transaction found in incoming batch as a duplicate, not marking as suspicious")
						break
					}
				}
			}

			if !found && len(dbt.SuspiciousReasons) == 0 {
				// Mark as suspicious
				// Mark as suspicious
				dbtNoIdFull := models.TransactionWithoutID(&dbt)
				dbtNoIdFull.SuspiciousReasons = []string{"Not present in importer transactions"}

				_, err := s.db.UpdateTransaction(userID, dbt.Id, dbtNoIdFull)
				if err != nil {
					// Ignore if not found (deleted)
					s.logger.With("error", err, "transactionID", dbt.Id).Warn("Failed to mark transaction as suspicious (might be deleted)")
				} else {
					suspiciousCnt++
				}
			}
		}
		if suspiciousCnt > 0 {
			_, err := s.db.CreateNotification(userID, &goserver.Notification{
				Date:        time.Now(),
				Type:        string(models.NotificationTypeInfo),
				Title:       "Suspicious Transactions Detected",
				Description: fmt.Sprintf("Import from %q found %d transactions that were not present in the bank data. Please review them.", biData.Name, suspiciousCnt),
			})
			if err != nil {
				s.logger.With("error", err).Error("Failed to create notification for suspicious transactions")
			}
		}
	}

	// if all transactions are processed then update account bank info
	if checkMissing && info != nil && len(info.Balances) > 0 {
		acc, err := s.db.GetAccount(userID, biData.AccountId)
		if err == nil {
			accNoId := models.AccountWithoutID(&acc)
			// Update account bank info from the importer provided info
			if info.AccountId != "" {
				accNoId.BankInfo.AccountId = info.AccountId
			}
			if info.BankId != "" {
				accNoId.BankInfo.BankId = info.BankId
			}

			// Update or append balances
			for _, b := range info.Balances {
				found := false
				for i := range accNoId.BankInfo.Balances {
					if accNoId.BankInfo.Balances[i].CurrencyId == b.CurrencyId {
						accNoId.BankInfo.Balances[i] = b
						found = true
						break
					}
				}
				if !found {
					accNoId.BankInfo.Balances = append(accNoId.BankInfo.Balances, b)
				}
			}

			_, err = s.db.UpdateAccount(userID, biData.AccountId, accNoId)
			if err != nil {
				s.logger.With("error", err).Error("Failed to update account bank info")
			} else {
				s.logger.Info("Updated account bank info", "accountId", biData.AccountId)
			}
		} else {
			s.logger.With("error", err, "accountId", biData.AccountId).Error("Failed to get account for bank info update")
		}
	}

	// update last import fields
	lastImport, err := s.updateLastImportFields(userID, id, info, len(transactions), cnt, suspiciousCnt)
	if err != nil {
		return nil, fmt.Errorf("can't update last import fields: %w", err)
	}

	return lastImport, nil
}

func (s *BankImportersAPIServiceImpl) Upload(
	userID, id, format string, data []byte, containsAllTransactions bool,
) (*goserver.ImportResult, error) {
	biData, err := s.db.GetBankImporter(userID, id)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get BankImporter")
		return nil, fmt.Errorf("can't get BankImporter: %w", err)
	}

	currencies, err := s.db.GetCurrencies(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get currencies")
		return nil, fmt.Errorf("can't get currencies: %w", err)
	}

	cp := bankimporters.NewDefaultCurrencyProvider(s.db, userID, currencies)

	var bi bankimporters.Importer

	switch biData.Type {
	case "revolut":
		bi, err = bankimporters.NewRevolutConverter(s.logger, biData, cp)
		if err != nil {
			s.logger.With("error", err).Error("Failed to create RevolutConverter")
			return nil, fmt.Errorf("can't create RevolutConverter: %w", err)
		}
	case "kb":
		bi, err = bankimporters.NewKBConverter(s.logger, biData, cp)
		if err != nil {
			s.logger.With("error", err).Error("Failed to create KbConverter")
			return nil, fmt.Errorf("can't create KbConverter: %w", err)
		}
	default:
		s.logger.With("type", biData.Type).Error("Unsupported bank importer type")
		_ = s.addImportResult(userID, id, goserver.ImportResult{
			Date:        time.Now(),
			Status:      "error",
			Description: fmt.Sprintf("Unsupported bank importer type: %s", biData.Type),
		})
		return nil, fmt.Errorf("unsupported bank importer type: %s", biData.Type)
	}

	info, transactions, err := bi.ParseAndImport(format, string(data))
	if err != nil {
		s.logger.With("error", err).Error("Failed to parse and import data")
		_ = s.addImportResult(userID, id, goserver.ImportResult{
			Date:        time.Now(),
			Status:      "error",
			Description: fmt.Sprintf("Failed to parse and import data: %s", err),
		})
		return nil, fmt.Errorf("can't parse and import data: %w", err)
	}

	lastImport, err := s.saveImportedTransactions(userID, id, info, transactions, containsAllTransactions)
	if err != nil {
		s.logger.With("error", err).Error("Failed to save imported transactions")
		return nil, fmt.Errorf("can't save imported transactions: %w", err)
	}

	return lastImport, nil
}

func (s *BankImportersAPIServiceImpl) UploadBankImporter(
	ctx context.Context, id, format string, containsAllTransactions bool, file *os.File,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	data, err := os.ReadFile(file.Name())
	if err != nil {
		s.logger.With("error", err).Error("Failed to read uploaded file")
		return goserver.Response(500, nil), nil
	}

	lastImport, err := s.Upload(userID, id, format, data, containsAllTransactions)
	if err != nil {
		s.logger.With("error", err).Error("Failed to upload")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, lastImport), nil
}

func isPerfectMatch(m *goserver.Matcher) bool {
	if len(m.ConfirmationHistory) < 10 {
		return false
	}
	for _, confirmed := range m.ConfirmationHistory {
		if !confirmed {
			return false
		}
	}
	return true
}

func getIncreases(movements []goserver.Movement) map[string]float64 {
	pos := make(map[string]float64)
	neg := make(map[string]float64)
	for _, m := range movements {
		if m.Amount > 0 {
			pos[m.CurrencyId] += m.Amount
		} else {
			neg[m.CurrencyId] -= m.Amount
		}
	}

	res := make(map[string]float64)
	for c, p := range pos {
		n := neg[c]
		if p > n {
			res[c] = p
		} else {
			res[c] = n
		}
	}
	for c, n := range neg {
		if _, ok := res[c]; !ok {
			res[c] = n
		}
	}
	return res
}

func isDuplicate(t1 *goserver.TransactionNoId, t2 *goserver.Transaction) bool {
	// 1. Time check (+/- 2 days)
	delta := t1.Date.Sub(t2.Date)
	if delta < 0 {
		delta = -delta
	}
	if delta > 2*time.Hour*24 {
		return false
	}

	// 2. Amount check (sum of increases per currency)
	inc1 := getIncreases(t1.Movements)
	inc2 := getIncreases(t2.Movements)

	if len(inc1) != len(inc2) {
		return false
	}

	for c, v1 := range inc1 {
		v2, ok := inc2[c]
		if !ok || math.Abs(v1-v2) >= 0.01 {
			return false
		}
	}

	return true
}
