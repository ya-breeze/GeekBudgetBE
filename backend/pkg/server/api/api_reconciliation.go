package api

import (
	"context"
	"log/slog"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/constants"
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
	familyID, ok := constants.GetFamilyID(ctx)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	accounts, err := s.db.GetAccounts(familyID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get accounts")
		return goserver.Response(500, nil), nil
	}

	bankImporters, err := s.db.GetBankImporters(familyID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get bank importers")
		return goserver.Response(500, nil), nil
	}

	currencies, err := s.db.GetCurrencies(familyID)
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

	bulkData, err := s.db.GetBulkReconciliationData(familyID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get bulk reconciliation data")
		return goserver.Response(500, nil), nil
	}

	var statuses []goserver.ReconciliationStatus
	for _, acc := range accounts {
		if acc.Type != "asset" {
			continue
		}

		for _, b := range acc.BankInfo.Balances {
			var lastRec *goserver.Reconciliation
			if accLatest, ok := bulkData.LatestReconciliations[acc.Id]; ok {
				lastRec = accLatest[b.CurrencyId]
			}

			hasImporter := accountsWithImporter[acc.Id]
			// Filter: only show if it has a bank importer OR it was already manually reconciled OR it's explicitly marked
			if !hasImporter && lastRec == nil && !acc.ShowInReconciliation {
				continue
			}

			appBalance := bulkData.Balances[acc.Id][b.CurrencyId]
			unprocessedCount := bulkData.UnprocessedCounts[acc.Id]

			bankBalance := b.ClosingBalance
			var bankBalanceAt *time.Time = b.LastUpdatedAt

			// If no importer but we have manual reconciliation, use that as "Bank Balance"
			if !hasImporter && lastRec != nil {
				bankBalance = lastRec.ReconciledBalance
				bankBalanceAt = &lastRec.ReconciledAt
			}

			status := goserver.ReconciliationStatus{
				AccountId:                  acc.Id,
				AccountName:                acc.Name,
				CurrencyId:                 b.CurrencyId,
				CurrencySymbol:             currencyMap[b.CurrencyId],
				BankBalance:                bankBalance,
				AppBalance:                 appBalance,
				Delta:                      appBalance.Sub(bankBalance),
				HasUnprocessedTransactions: unprocessedCount > 0,
				HasBankImporter:            hasImporter,
				BankBalanceAt:              bankBalanceAt,
			}

			if bankBalanceAt != nil {
				maxDate := bulkData.MaxTransactionDates[acc.Id][b.CurrencyId]
				status.HasTransactionsAfterBankBalance = maxDate.After(*bankBalanceAt)
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
	familyID, ok := constants.GetFamilyID(ctx)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Detect whether this account has a bank importer configured.
	// Use GetBankImporters (same as GetReconciliationStatus) rather than inspecting
	// BankInfo.Balances — an importer that has never run would have no Balances entries
	// but is still a configured importer and should enforce the tolerance check.
	importers, err := s.db.GetBankImporters(familyID)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get bank importers")
		return goserver.Response(500, nil), nil
	}
	hasImporter := false
	for _, imp := range importers {
		if imp.AccountId == id {
			hasImporter = true
			break
		}
	}

	// Resolve balance: if frontend sends 0, use the current computed account balance.
	balance := body.Balance
	if balance.IsZero() {
		balance, err = s.db.GetAccountBalance(familyID, id, body.CurrencyId)
		if err != nil {
			s.logger.With("error", err).Error("Failed to get account balance")
			return goserver.Response(500, nil), nil
		}
	}

	var expectedBalance decimal.Decimal
	isManual := true

	if hasImporter {
		// Importer path: derive expected balance from last import data and enforce tolerance.
		acc, accErr := s.db.GetAccount(familyID, id)
		if accErr != nil {
			return goserver.Response(404, nil), nil
		}
		for _, b := range acc.BankInfo.Balances {
			if b.CurrencyId == body.CurrencyId {
				expectedBalance = b.ClosingBalance
				break
			}
		}
		if balance.Sub(expectedBalance).Abs().GreaterThan(constants.ReconciliationTolerance) {
			return goserver.Response(400, "Cannot reconcile: account balance does not match bank balance"), nil
		}
		isManual = body.Balance.IsPositive()
	} else {
		// No-importer path: the user is confirming the app balance is correct.
		// Set expectedBalance = balance so the history record shows delta = 0.
		// IsManual is always true for no-importer accounts.
		expectedBalance = balance
	}

	rec, err := s.db.CreateReconciliation(familyID, &goserver.ReconciliationNoId{
		AccountId:         id,
		CurrencyId:        body.CurrencyId,
		ReconciledBalance: balance,
		ExpectedBalance:   expectedBalance,
		IsManual:          isManual,
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
	familyID, ok := constants.GetFamilyID(ctx)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Get last reconciliation from new entity
	lastRec, err := s.db.GetLatestReconciliation(familyID, id, currencyId)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get reconciliation")
		return goserver.Response(500, nil), nil
	}

	var dateFrom time.Time
	if lastRec != nil {
		dateFrom = lastRec.ReconciledAt
	}

	// Get all transactions from that date
	allTransactions, err := s.db.GetTransactions(familyID, dateFrom, time.Time{}, false)
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
	familyID, ok := constants.GetFamilyID(ctx)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Create initial reconciliation record
	rec, err := s.db.CreateReconciliation(familyID, &goserver.ReconciliationNoId{
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

// GetReconciliationHistory returns all reconciliation records for an account+currency pair
func (s *ReconciliationAPIServiceImpl) GetReconciliationHistory(ctx context.Context, id string, currencyId string) (goserver.ImplResponse, error) {
	familyID, ok := constants.GetFamilyID(ctx)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	recs, err := s.db.GetReconciliationsForAccountAndCurrency(familyID, id, currencyId)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get reconciliation history")
		return goserver.Response(500, nil), nil
	}

	return goserver.Response(200, recs), nil
}

// AnalyzeDisbalance find transactions that might explain the disbalance
func (s *ReconciliationAPIServiceImpl) AnalyzeDisbalance(ctx context.Context, id string, body goserver.AnalyzeDisbalanceRequest) (goserver.ImplResponse, error) {
	familyID, ok := constants.GetFamilyID(ctx)
	if !ok {
		return goserver.Response(500, nil), nil
	}

	// Fetch transactions since last reconciliation
	lastRec, err := s.db.GetLatestReconciliation(familyID, id, body.CurrencyId)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get reconciliation")
		return goserver.Response(500, nil), nil
	}

	var dateFrom time.Time
	if lastRec != nil {
		dateFrom = lastRec.ReconciledAt
	}

	allTransactions, err := s.db.GetTransactions(familyID, dateFrom, time.Time{}, false)
	if err != nil {
		s.logger.With("error", err).Error("Failed to get transactions")
		return goserver.Response(500, nil), nil
	}

	// Filter to only transactions affecting this account+currency
	var filtered []goserver.Transaction
	for _, tx := range allTransactions {
		for _, m := range tx.Movements {
			if m.AccountId == id && m.CurrencyId == body.CurrencyId {
				filtered = append(filtered, tx)
				break
			}
		}
	}

	analysis := common.AnalyzeDisbalance(body.TargetDelta, filtered, id, body.CurrencyId)

	return goserver.Response(200, analysis), nil
}
