package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/bankimporters"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
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
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	res, err := s.db.CreateBankImporter(userID, &input)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create BankImporter")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, res), nil
}

func (s *BankImportersAPIServiceImpl) DeleteBankImporter(ctx context.Context, id string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
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
	userID, ok := ctx.Value(UserIDKey).(string)
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
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	res, err := s.db.UpdateBankImporter(userID, id, &input)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update BankImporter")
		return goserver.Response(500, nil), nil
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
	userID, ok := ctx.Value(UserIDKey).(string)
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

func (s *BankImportersAPIServiceImpl) fetchFioTransactions(ctx context.Context, userID, id string,
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

	bi, err := bankimporters.NewFioConverter(s.logger, biData, currencies)
	if err != nil {
		return nil, nil, fmt.Errorf("can't create FioConverter: %w", err)
	}

	s.logger.Info("Importing transactions")
	info, transactions, err := bi.Import(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("can't import transactions: %w", err)
	}
	s.logger.With("info", info, "transactions", len(transactions)).Info("Imported transactions")

	return info, transactions, nil
}

func (s *BankImportersAPIServiceImpl) updateLastImportFields(
	userID string, id string, info *goserver.BankAccountInfo, totalTransactionsCnt int, newTransactionsCnt int,
) (*goserver.ImportResult, error) {
	biData, err := s.db.GetBankImporter(userID, id)
	if err != nil {
		return nil, fmt.Errorf("can't fetch bank importer: %w", err)
	}
	lastImport := goserver.ImportResult{
		Date:   biData.LastSuccessfulImport,
		Status: "success",
		Description: fmt.Sprintf("Fetched %d transactions. Imported %d new transactions. Final balances: %v",
			totalTransactionsCnt, newTransactionsCnt, info.Balances),
	}
	biData.LastSuccessfulImport = time.Now()
	biData.LastImports = append(biData.LastImports, lastImport)
	if len(biData.LastImports) > 10 {
		biData.LastImports = biData.LastImports[1:]
	}
	_, err = s.db.UpdateBankImporter(userID, id, &biData)
	if err != nil {
		return nil, fmt.Errorf("can't update BankImporter: %w", err)
	}

	return &lastImport, nil
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

	// save transactions to the database
	cnt := 0
	for _, t := range transactions {
		if len(t.ExternalIds) != 1 {
			return nil, fmt.Errorf("transaction has invalid external IDs: %v", t)
		}

		// search for existing transaction with the same external ID. If found, skip saving
		found := false
		for _, dbt := range dbTransactions {
			if slices.Contains(dbt.ExternalIds, t.ExternalIds[0]) {
				found = true
				s.logger.With("externalID", t.ExternalIds[0]).Info("Transaction already was imported")
				break
			}
		}
		if found {
			continue
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

func (s *BankImportersAPIServiceImpl) UploadBankImporter(
	ctx context.Context, id, format string, file *os.File,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	biData, err := s.db.GetBankImporter(userID, id)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get BankImporter")
		return goserver.Response(500, nil), nil
	}

	if biData.Type != "revolut" {
		s.logger.With("type", biData.Type).Error("Unsupported bank importer type")
		return goserver.Response(500, nil), nil
	}

	currencies, err := s.db.GetCurrencies(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get currencies")
		return goserver.Response(500, nil), nil
	}

	bi, err := bankimporters.NewRevolutConverter(s.logger, biData, currencies)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create RevolutConverter")
		return goserver.Response(500, nil), nil
	}

	data, err := os.ReadFile(file.Name())
	if err != nil {
		s.logger.With("error", err).Error("Failed to read uploaded file")
		return goserver.Response(500, nil), nil
	}

	info, transactions, err := bi.ParseAndImport(ctx, format, string(data))
	if err != nil {
		s.logger.With("error", err).Error("Failed to parse and import Revolut data")
		return goserver.Response(500, nil), nil
	}

	lastImport, err := s.saveImportedTransactions(userID, id, info, transactions)
	if err != nil {
		s.logger.With("error", err).Error("Failed to save imported transactions")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, lastImport), nil
}
