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

func (s *BankImportersAPIServiceImpl) CreateBankImporter(context.Context, goserver.BankImporterNoId,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *BankImportersAPIServiceImpl) DeleteBankImporter(context.Context, string) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *BankImportersAPIServiceImpl) GetBankImporters(context.Context) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *BankImportersAPIServiceImpl) UpdateBankImporter(context.Context, string, goserver.BankImporterNoId,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}
