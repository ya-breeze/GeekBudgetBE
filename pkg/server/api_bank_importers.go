package server

import (
	"context"
	"log/slog"
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

func NewBankImportersAPIServiceImpl(logger *slog.Logger, db database.Storage,
) goserver.BankImportersAPIServicer {
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

func (s *BankImportersAPIServiceImpl) FetchBankImporter(
	ctx context.Context, id string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}
	s.logger.With("user", userID).Info("Fetching transactions for bank importer")

	biData, err := s.db.GetBankImporter(userID, id)
	if err != nil {
		s.logger.With("error", err).Error("Failed to fetch for bank importer")
		return goserver.Response(500, nil), nil
	}

	bi, err := bankimporters.NewFioConverter(s.logger, biData)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create FioConverter")
		return goserver.Response(500, nil), nil
	}

	s.logger.Info("Importing transactions")
	info, transactions, err := bi.Import(ctx)
	if err != nil {
		s.logger.With("error", err).Error("Failed to import transactions")
		return goserver.Response(500, nil), nil
	}
	s.logger.With("info", info, "transactions", len(transactions)).Info("Imported transactions")

	// Fetch all transactions from the database
	// TODO don't fetch - just search by external ID
	dbTransactions, err := s.db.GetTransactions(userID, time.Time{}, time.Time{})
	if err != nil {
		s.logger.With("error", err).Error("Failed to fetch transactions")
		return goserver.Response(500, nil), nil
	}

	// save transactions to the database
	for _, t := range transactions {
		if len(t.ExternalIds) != 1 {
			s.logger.With("transaction", t).Error("Transaction has invalid external IDs")
			return goserver.Response(500, nil), nil
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
			s.logger.With("error", err).Error("Failed to save transaction")
			return goserver.Response(500, nil), nil
		}
		s.logger.Info("Imported transaction saved to DB")
	}
	s.logger.Info("All imported transactions saved to DB")

	// update last import fields
	biData.LastSuccessfulImport = time.Now()
	biData.LastImports = append(biData.LastImports, goserver.BankImporterNoIdLastImportsInner{
		Date:   biData.LastSuccessfulImport,
		Status: "OK",
	})
	_, err = s.db.UpdateBankImporter(userID, id, &biData)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update BankImporter")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, nil), nil
}
