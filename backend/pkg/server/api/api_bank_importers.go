package api

import (
	"context"
	"fmt"
	"log/slog"
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
	ctx context.Context, userID, importerID string,
) (*goserver.ImportResult, error) {
	s.logger.Info("Fetching bank importer", "userID", userID, "bankImporterID", importerID)
	info, transactions, err := s.fetchFioTransactions(ctx, userID, importerID)
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

	lastImport, err := s.saveImportedTransactions(userID, importerID, info, transactions)
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

	lastImport, err := s.Fetch(ctx, userID, id)
	if err != nil {
		s.logger.With("error", err).Error("Failed to fetch")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, lastImport), nil
}

func (s *BankImportersAPIServiceImpl) fetchFioTransactions(
	ctx context.Context, userID, id string,
) (*goserver.BankAccountInfo, []goserver.TransactionNoId, error) {
	s.logger.With("user", userID).Info("Fetching transactions for bank importer")

	biData, err := s.db.GetBankImporter(userID, id)
	if err != nil {
		return nil, nil, fmt.Errorf("can't fetch bank importer: %w", err)
	}

	currencies, err := s.db.GetCurrencies(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("can't fetch currencies: %w", err)
	}

	cp := bankimporters.NewDefaultCurrencyProvider(s.db, userID, currencies)

	bi, err := bankimporters.NewFioConverter(s.logger, biData, cp)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create FioConverter: %w", err)
	}

	s.logger.Info("Importing transactions")
	info, transactions, err := bi.Import(ctx)
	if err != nil {
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
		}

		return nil, nil, fmt.Errorf("can't import transactions: %w", err)
	}
	s.logger.With("info", info, "transactions", len(transactions)).Info("Imported transactions")

	if biData.FetchAll {
		s.logger.With("bankImporterID", id).With("userID", userID).Info("All transactions fetched. Disabling FetchAll")
		biData.FetchAll = false
		_, err = s.db.UpdateBankImporter(userID, id, &biData)
		if err != nil {
			return nil, nil, fmt.Errorf("can't update BankImporter: %w", err)
		}
	}

	return info, transactions, nil
}

func (s *BankImportersAPIServiceImpl) updateLastImportFields(
	userID string, id string, info *goserver.BankAccountInfo, totalTransactionsCnt int, newTransactionsCnt int,
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
		Balances: balances,
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
	userID, id string, info *goserver.BankAccountInfo, transactions []goserver.TransactionNoId,
) (*goserver.ImportResult, error) {
	// Fetch all transactions from the database
	// TODO don't fetch - just search by external ID
	dbTransactions, err := s.db.GetTransactions(userID, time.Time{}, time.Time{})
	if err != nil {
		return nil, fmt.Errorf("can't fetch transactions from DB: %w", err)
	}

	matchers, err := s.db.GetMatchersRuntime(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matchers")
		return nil, fmt.Errorf("can't get matchers: %w", err)
	}

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
		tStableHash := bankimporters.ComputeStableHash(&t)

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

			// 2. Check for stable hash match
			// We need to construct a TransactionNoId from dbt to compute the hash
			dbtNoId := goserver.TransactionNoId{
				Date:      dbt.Date,
				Movements: dbt.Movements,
			}
			dbtStableHash := bankimporters.ComputeStableHash(&dbtNoId)
			if dbtStableHash == tStableHash {
				found = true
				s.logger.With("stableHash", tStableHash).Info("Transaction already was imported (stable hash match)")

				// Optional: Update the existing transaction with the new legacy hash?
				// If we matched by stable hash, it means the legacy hash is missing from dbt.
				// We probably should add it to dbt.ExternalIds to prevent future re-computations?
				// But mutating dbt here is complex (need to save to DB).
				// For now, accept it as duplicate and skip.
				break
			}
		}
		if found {
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

		for _, matcher := range matchers {
			if common.Match(&matcher, tempDetails) != common.MatchResultSuccess {
				continue
			}

			if isPerfectMatch(matcher.Matcher) {
				s.logger.Info("Found perfect match", "matcher", matcher.Matcher.OutputDescription, "transaction", t.Description)

				t.Description = matcher.Matcher.OutputDescription
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
				break
			}
		}

		_, err = s.db.CreateTransaction(userID, &t)
		if err != nil {
			return nil, fmt.Errorf("can't save transaction: %w", err)
		}
		s.logger.Info("Imported transaction saved to DB")
		cnt++
	}
	s.logger.With("count", cnt).Info("All new imported transactions saved to DB")

	// update last import fields
	lastImport, err := s.updateLastImportFields(userID, id, info, len(transactions), cnt)
	if err != nil {
		return nil, fmt.Errorf("can't update last import fields: %w", err)
	}

	return lastImport, nil
}

func (s *BankImportersAPIServiceImpl) Upload(
	userID, id, format string, data []byte,
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

	lastImport, err := s.saveImportedTransactions(userID, id, info, transactions)
	if err != nil {
		s.logger.With("error", err).Error("Failed to save imported transactions")
		return nil, fmt.Errorf("can't save imported transactions: %w", err)
	}

	return lastImport, nil
}

func (s *BankImportersAPIServiceImpl) UploadBankImporter(
	ctx context.Context, id, format string, file *os.File,
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

	lastImport, err := s.Upload(userID, id, format, data)
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
