package server

import (
	"context"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type CurrenciesAPIServicerImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewCurrenciesAPIServicer(logger *slog.Logger, db database.Storage) goserver.CurrenciesAPIServicer {
	return &CurrenciesAPIServicerImpl{
		logger: logger,
		db:     db,
	}
}

func (s *CurrenciesAPIServicerImpl) GetCurrencies(ctx context.Context) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	currencies, err := s.db.GetCurrencies(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get currencies")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, currencies), nil
}

func (s *CurrenciesAPIServicerImpl) CreateCurrency(
	ctx context.Context, currencyNoId goserver.CurrencyNoId) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	currency, err := s.db.CreateCurrency(userID, &currencyNoId)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create currency")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, currency), nil
}

func (s *CurrenciesAPIServicerImpl) UpdateCurrency(
	ctx context.Context, currencyID string, currencyNoID goserver.CurrencyNoId) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	currency, err := s.db.UpdateCurrency(userID, currencyID, &currencyNoID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update currency")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, currency), nil
}

func (s *CurrenciesAPIServicerImpl) DeleteCurrency(
	ctx context.Context, currencyID string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	err := s.db.DeleteCurrency(userID, currencyID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to delete currency")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, nil), nil
}

func (s *CurrenciesAPIServicerImpl) GetCurrency(
	ctx context.Context, currencyID string) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	currency, err := s.db.GetCurrency(userID, currencyID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get currency")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, currency), nil
}
