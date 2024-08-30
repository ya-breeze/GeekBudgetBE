package server

import (
	"context"
	"log/slog"

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
	_, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// err := s.db.FetchBankImporter(userID, id)
	// if err != nil {
	// s.logger.With("error", err).Error("Failed to fetch for bank importer")
	// return goserver.Response(500, nil), nil
	// }

	// return goserver.Response(200, nil), nil
	return goserver.Response(500, "not implemented yet"), nil
}
