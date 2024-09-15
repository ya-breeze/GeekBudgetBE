package server

import (
	"context"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type ImportAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewImportAPIServiceImpl(logger *slog.Logger, db database.Storage,
) goserver.ImportAPIServicer {
	return &ImportAPIServiceImpl{logger: logger, db: db}
}

//nolint:cyclop // This function is not complex
func (s *ImportAPIServiceImpl) CallImport(ctx context.Context, data goserver.WholeUserData,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	for _, currency := range data.Currencies {
		if _, err := s.db.UpdateCurrency(userID, currency.Id, models.CurrencyWithoutID(&currency)); err != nil {
			s.logger.With("error", err).Error("Failed to add currency")
			return goserver.Response(500, nil), nil
		}
	}

	for _, account := range data.Accounts {
		if _, err := s.db.UpdateAccount(userID, account.Id, models.AccountWithoutID(&account)); err != nil {
			s.logger.With("error", err).Error("Failed to add account")
			return goserver.Response(500, nil), nil
		}
	}

	for _, transaction := range data.Transactions {
		if _, err := s.db.UpdateTransaction(userID, transaction.Id, models.TransactionWithoutID(&transaction)); err != nil {
			s.logger.With("error", err).Error("Failed to add transaction")
			return goserver.Response(500, nil), nil
		}
	}

	for _, matcher := range data.Matchers {
		if _, err := s.db.UpdateMatcher(userID, matcher.Id, models.MatcherWithoutID(&matcher)); err != nil {
			s.logger.With("error", err).Error("Failed to add matcher")
			return goserver.Response(500, nil), nil
		}
	}

	for _, bankImporter := range data.BankImporters {
		if _, err := s.db.UpdateBankImporter(userID, bankImporter.Id,
			models.BankImporterWithoutID(&bankImporter)); err != nil {
			s.logger.With("error", err).Error("Failed to add bank importer")
			return goserver.Response(500, nil), nil
		}
	}

	return goserver.Response(200, nil), nil
}
