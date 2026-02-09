package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/database"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"github.com/ya-breeze/geekbudgetbe/pkg/server/common"
)

type ReconciliationAPIServiceImpl struct {
	logger *slog.Logger
	db     database.Storage
}

func NewReconciliationAPIServiceImpl(logger *slog.Logger, db database.Storage) *ReconciliationAPIServiceImpl {
	return &ReconciliationAPIServiceImpl{logger: logger, db: db}
}

// GetReconciliationStatus returns reconciliation status for all asset accounts
func (s *ReconciliationAPIServiceImpl) GetReconciliationStatus(ctx context.Context) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	accounts, err := s.db.GetAccounts(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts")
		return goserver.Response(500, nil), nil
	}

	bankImporters, err := s.db.GetBankImporters(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get bank importers")
		return goserver.Response(500, nil), nil
	}

	currencies, err := s.db.GetCurrencies(userID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get currencies")
		return goserver.Response(500, nil), nil
	}

	currencyMap := make(map[string]string)
	for _, c := range currencies {
		currencyMap[c.Id] = c.Name
	}

	// Build map of accounts with bank importers
	accountsWithImporter := make(map[string]bool)
	for _, bi := range bankImporters {
		accountsWithImporter[bi.AccountId] = true
	}

	var statuses []goserver.ReconciliationStatus
	for _, acc := range accounts {
		if acc.Type != "asset" {
			continue
		}

		for _, b := range acc.BankInfo.Balances {
			// Get latest reconciliation from new entity
			lastRec, _ := s.db.GetLatestReconciliation(userID, acc.Id, b.CurrencyId)

			hasImporter := accountsWithImporter[acc.Id]
			// Filter: only show if it has a bank importer OR it was already manually reconciled
			if !hasImporter && lastRec == nil {
				continue
			}

			appBalance, err := s.db.GetAccountBalance(userID, acc.Id, b.CurrencyId)
			if err != nil {
				s.logger.With("error", err, "accountId", acc.Id, "currencyId", b.CurrencyId).Warn("Failed to get account balance")
				continue
			}

			unprocessedCount, _ := s.db.CountUnprocessedTransactionsForAccount(userID, acc.Id, acc.IgnoreUnprocessedBefore)

			status := goserver.ReconciliationStatus{
				AccountId:                  acc.Id,
				AccountName:                acc.Name,
				CurrencyId:                 b.CurrencyId,
				CurrencySymbol:             currencyMap[b.CurrencyId],
				BankBalance:                b.ClosingBalance,
				AppBalance:                 appBalance,
				Delta:                      appBalance - b.ClosingBalance,
				HasUnprocessedTransactions: unprocessedCount > 0,
				HasBankImporter:            hasImporter,
			}

			if lastRec != nil {
				status.LastReconciledAt = &lastRec.ReconciledAt
				status.LastReconciledBalance = lastRec.ReconciledBalance
				status.IsManualReconciliationEnabled = lastRec.IsManual
			}

			statuses = append(statuses, status)
		}
	}

	return goserver.Response(200, statuses), nil
}

// ReconcileAccount creates a new reconciliation record
func (s *ReconciliationAPIServiceImpl) ReconcileAccount(
	ctx context.Context, id string, body goserver.ReconcileAccountRequest,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Get current app balance if not provided
	balance := body.Balance
	if balance == 0 {
		var err error
		balance, err = s.db.GetAccountBalance(userID, id, body.CurrencyId)
		if err != nil {
			s.logger.With("error", err).Error("Failed to get account balance")
			return goserver.Response(500, nil), nil
		}
	}

	// Get bank balance for expected balance
	acc, err := s.db.GetAccount(userID, id)
	if err != nil {
		return goserver.Response(404, nil), nil
	}

	var expectedBalance float64
	for _, b := range acc.BankInfo.Balances {
		if b.CurrencyId == body.CurrencyId {
			expectedBalance = b.ClosingBalance
			break
		}
	}

	// Create new reconciliation record
	rec, err := s.db.CreateReconciliation(userID, &goserver.ReconciliationNoId{
		AccountId:         id,
		CurrencyId:        body.CurrencyId,
		ReconciledBalance: balance,
		ExpectedBalance:   expectedBalance,
		IsManual:          body.Balance > 0, // Manual if balance explicitly provided
	})
	if err != nil {
		s.logger.With("error", err).Error("Failed to create reconciliation")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, rec), nil
}

// GetTransactionsSinceReconciliation returns transactions since last reconciliation
func (s *ReconciliationAPIServiceImpl) GetTransactionsSinceReconciliation(
	ctx context.Context, id string, currencyId string,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Get last reconciliation from new entity
	lastRec, err := s.db.GetLatestReconciliation(userID, id, currencyId)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get reconciliation")
		return goserver.Response(500, nil), nil
	}

	var dateFrom time.Time
	if lastRec != nil {
		dateFrom = lastRec.ReconciledAt
	}

	// Get all transactions from that date
	allTransactions, err := s.db.GetTransactions(userID, dateFrom, time.Time{}, false)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions")
		return goserver.Response(500, nil), nil
	}

	// Filter to only transactions affecting this account+currency
	var filtered []goserver.Transaction
	for _, tx := range allTransactions {
		for _, m := range tx.Movements {
			if m.AccountId == id && m.CurrencyId == currencyId {
				filtered = append(filtered, tx)
				break
			}
		}
	}

	return goserver.Response(200, filtered), nil
}

// EnableAccountReconciliation creates initial reconciliation for manual accounts
func (s *ReconciliationAPIServiceImpl) EnableAccountReconciliation(
	ctx context.Context, id string, body goserver.EnableReconciliationRequest,
) (goserver.ImplResponse, error) {
	userID, ok := ctx.Value(common.UserIDKey).(string)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Create initial reconciliation record
	rec, err := s.db.CreateReconciliation(userID, &goserver.ReconciliationNoId{
		AccountId:         id,
		CurrencyId:        body.CurrencyId,
		ReconciledBalance: body.InitialBalance,
		ExpectedBalance:   body.InitialBalance,
		IsManual:          true,
	})
	if err != nil {
		s.logger.With("error", err).Error("Failed to create reconciliation")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, rec), nil
}
