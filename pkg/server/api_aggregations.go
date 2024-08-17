package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type AggregationsAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewAggregationsAPIServiceImpl(logger *slog.Logger, db database.Storage,
) goserver.AggregationsAPIServicer {
	return &AggregationsAPIServiceImpl{logger: logger, db: db}
}

func (s *AggregationsAPIServiceImpl) GetBalances(context.Context, time.Time, time.Time, string,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *AggregationsAPIServiceImpl) GetExpenses(context.Context, time.Time, time.Time, string,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}

func (s *AggregationsAPIServiceImpl) GetIncomes(context.Context, time.Time, time.Time, string,
) (goserver.ImplResponse, error) {
	return goserver.ImplResponse{}, nil
}
