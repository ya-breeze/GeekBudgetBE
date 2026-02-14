package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type ExportAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewExportAPIServiceImpl(logger *slog.Logger, db database.Storage,
) goserver.ExportAPIServicer {
	return &ExportAPIServiceImpl{logger: logger, db: db}
}

func (s *ExportAPIServiceImpl) Export(ctx context.Context) (goserver.ImplResponse, error) {
	var err error
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	res := goserver.WholeUserData{}

	res.Currencies, err = s.db.GetCurrencies(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get currencies")
		return goserver.Response(500, nil), nil
	}

	res.Accounts, err = s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts")
		return goserver.Response(500, nil), nil
	}

	res.Transactions, err = s.db.GetTransactions(userID, time.Time{}, time.Time{}, false)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions")
		return goserver.Response(500, nil), nil
	}

	res.Matchers, err = s.db.GetMatchers(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get matchers")
		return goserver.Response(500, nil), nil
	}

	res.BankImporters, err = s.db.GetBankImporters(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get bank importers")
		return goserver.Response(500, nil), nil
	}

	user, err := s.db.GetUser(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get user")
		return goserver.Response(500, nil), nil
	}
	res.User = user.FromDB()

	return goserver.Response(200, res), nil
}
