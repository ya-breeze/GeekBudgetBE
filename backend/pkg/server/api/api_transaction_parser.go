package api

import (
	"context"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

func (s *TransactionsAPIServiceImpl) ParseTransaction(
	ctx context.Context,
	req goserver.TransactionParseRequest,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(constants.UserIDKey).(string)
	if !ok {
		s.logger.Error("UserID not found in context")
		return goserver.Response(500, nil), nil
	}

	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts for parser")
		return goserver.Response(500, nil), nil
	}

	currencies, err := s.db.GetCurrencies(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get currencies for parser")
		return goserver.Response(500, nil), nil
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)
	transaction, warnings := common.ParseTransactionText(req.Text, accounts, currencies, today)

	return goserver.Response(200, goserver.TransactionParseResponse{
		Transaction: transaction,
		Warnings:    warnings,
	}), nil
}
