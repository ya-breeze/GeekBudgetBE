package server

import (
	"context"
	"log/slog"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type AccountsAPIServicerImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewAccountsAPIService(logger *slog.Logger, db database.Storage) goserver.AccountsAPIServicer {
	return &AccountsAPIServicerImpl{
		logger: logger,
		db:     db,
	}
}

func (s *AccountsAPIServicerImpl) CreateAccount(
	ctx context.Context, acc goserver.AccountNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	account, err := s.db.CreateAccount(userID, &acc)
	if err != nil {
		s.logger.With("error", err).Error("Failed to create account")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, account), nil
}

func (s *AccountsAPIServicerImpl) GetAccounts(
	ctx context.Context,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts")
		return goserver.Response(500, nil), nil
	}

	accounts = append(accounts, goserver.Account{
		Id:   "",
		Name: "Unknown account",
	})

	return goserver.Response(200, accounts), nil
}

func (s *AccountsAPIServicerImpl) UpdateAccount(
	ctx context.Context, accountID string, acc goserver.AccountNoId,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	account, err := s.db.UpdateAccount(userID, accountID, &acc)
	if err != nil {
		s.logger.With("error", err).Error("Failed to update account")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, account), nil
}

func (s *AccountsAPIServicerImpl) DeleteAccount(
	ctx context.Context, accountID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	if err := s.db.DeleteAccount(userID, accountID); err != nil {
		s.logger.With("error", err).Error("Failed to delete account")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, nil), nil
}

func (s *AccountsAPIServicerImpl) GetAccountHistory(
	ctx context.Context, accountID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	history, err := s.db.GetAccountHistory(userID, accountID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get account history")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, history), nil
}

func (s *AccountsAPIServicerImpl) GetAccount(
	ctx context.Context, accountID string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	account, err := s.db.GetAccount(userID, accountID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get account")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, account), nil
}
