package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

func (s *storage) GetLatestReconciliation(userID, accountID, currencyID string) (*goserver.Reconciliation, error) {
	var rec models.Reconciliation
	result := s.db.Where("user_id = ? AND account_id = ? AND currency_id = ?", userID, accountID, currencyID).
		Order("reconciled_at DESC").First(&rec)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // No reconciliation yet
		}
		return nil, fmt.Errorf("failed to get latest reconciliation: %w", result.Error)
	}
	return models.ReconciliationToAPI(&rec), nil
}

func (s *storage) GetReconciliationsForAccount(userID, accountID string) ([]goserver.Reconciliation, error) {
	var recs []models.Reconciliation
	err := s.db.Where("user_id = ? AND account_id = ?", userID, accountID).
		Order("reconciled_at DESC").Find(&recs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reconciliations for account: %w", err)
	}

	result := make([]goserver.Reconciliation, len(recs))
	for i := range recs {
		result[i] = *models.ReconciliationToAPI(&recs[i])
	}
	return result, nil
}

func (s *storage) GetReconciliationsForAccountAndCurrency(userID, accountID, currencyID string) ([]goserver.Reconciliation, error) {
	var recs []models.Reconciliation
	err := s.db.Where("user_id = ? AND account_id = ? AND currency_id = ?", userID, accountID, currencyID).
		Order("reconciled_at DESC").Find(&recs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get reconciliations for account and currency: %w", err)
	}

	result := make([]goserver.Reconciliation, len(recs))
	for i := range recs {
		result[i] = *models.ReconciliationToAPI(&recs[i])
	}
	return result, nil
}

func (s *storage) CreateReconciliation(userID string, rec *goserver.ReconciliationNoId) (goserver.Reconciliation, error) {
	model := models.ReconciliationFromAPI(userID, rec)
	model.ID = uuid.New()
	model.ReconciledAt = time.Now()

	if err := s.db.Create(&model).Error; err != nil {
		return goserver.Reconciliation{}, fmt.Errorf("failed to create reconciliation: %w", err)
	}

	if err := s.recordAuditLog(s.db, userID, "Reconciliation", model.ID.String(), "CREATED", &model); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	return *models.ReconciliationToAPI(model), nil
}

func (s *storage) InvalidateReconciliation(userID, accountID, currencyID string, fromDate time.Time) error {
	query := s.db.Where("user_id = ? AND account_id = ? AND currency_id = ?",
		userID, accountID, currencyID)

	if !fromDate.IsZero() {
		query = query.Where("reconciled_at >= ?", fromDate)
	}

	// Fetch reconciliation(s) to be deleted to record them in audit log
	var recs []models.Reconciliation
	if err := query.Find(&recs).Error; err != nil {
		return fmt.Errorf("failed to find reconciliations to invalidate: %w", err)
	}

	if len(recs) == 0 {
		return nil
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, rec := range recs {
			if err := s.recordAuditLog(tx, userID, "Reconciliation", rec.ID.String(), "DELETED", &rec); err != nil {
				s.log.Error("Failed to record audit log", "error", err)
			}
		}

		if err := tx.Delete(&recs).Error; err != nil {
			return fmt.Errorf("failed to invalidate reconciliations: %w", err)
		}
		return nil
	})
}

func (s *storage) GetBulkReconciliationData(userID string) (*BulkReconciliationData, error) {
	data := &BulkReconciliationData{
		Balances:              make(map[string]map[string]decimal.Decimal),
		LatestReconciliations: make(map[string]map[string]*goserver.Reconciliation),
		UnprocessedCounts:     make(map[string]int),
		MaxTransactionDates:   make(map[string]map[string]time.Time),
	}

	// 1. Get Accounts for opening balances and ignore dates
	accounts, err := s.GetAccounts(userID)
	if err != nil {
		return nil, err
	}
	ignoreMap := make(map[string]time.Time)
	for _, acc := range accounts {
		data.Balances[acc.Id] = make(map[string]decimal.Decimal)
		data.MaxTransactionDates[acc.Id] = make(map[string]time.Time)
		for _, b := range acc.BankInfo.Balances {
			data.Balances[acc.Id][b.CurrencyId] = b.OpeningBalance
		}
		if !acc.IgnoreUnprocessedBefore.IsZero() {
			ignoreMap[acc.Id] = acc.IgnoreUnprocessedBefore
		}
	}

	// 2. Get latest reconciliations
	var recs []models.Reconciliation
	err = s.db.Where("user_id = ?", userID).Order("reconciled_at DESC").Find(&recs).Error
	if err != nil {
		return nil, err
	}
	for i := range recs {
		r := &recs[i]
		apiRec := models.ReconciliationToAPI(r)
		if _, ok := data.LatestReconciliations[apiRec.AccountId]; !ok {
			data.LatestReconciliations[apiRec.AccountId] = make(map[string]*goserver.Reconciliation)
		}
		if _, ok := data.LatestReconciliations[apiRec.AccountId][apiRec.CurrencyId]; !ok {
			data.LatestReconciliations[apiRec.AccountId][apiRec.CurrencyId] = apiRec
		}
	}

	// 3. Get all transactions (merged_into_id IS NULL)
	var transactions []models.Transaction
	err = s.db.Where("user_id = ? AND merged_into_id IS NULL", userID).Find(&transactions).Error
	if err != nil {
		return nil, err
	}

	for _, t := range transactions {
		hasEmpty := false
		involvedAccounts := make(map[string]bool)

		for _, m := range t.Movements {
			// Update balances
			if m.AccountId != "" {
				if _, ok := data.Balances[m.AccountId]; ok {
					data.Balances[m.AccountId][m.CurrencyId] = data.Balances[m.AccountId][m.CurrencyId].Add(m.Amount)
				}
				involvedAccounts[m.AccountId] = true

				// Update MaxTransactionDates
				if _, ok := data.MaxTransactionDates[m.AccountId]; !ok {
					data.MaxTransactionDates[m.AccountId] = make(map[string]time.Time)
				}
				if t.Date.After(data.MaxTransactionDates[m.AccountId][m.CurrencyId]) {
					data.MaxTransactionDates[m.AccountId][m.CurrencyId] = t.Date
				}
			}

			if m.Amount.IsZero() {
				continue
			}
			if m.AccountId == "" {
				hasEmpty = true
			}
		}

		// Update unprocessed counts
		if hasEmpty {
			for accID := range involvedAccounts {
				ignoreDate := ignoreMap[accID]
				if ignoreDate.IsZero() || !t.Date.Before(ignoreDate) {
					data.UnprocessedCounts[accID]++
				}
			}
		}
	}

	return data, nil
}

func (s *storage) CountUnprocessedTransactionsForAccount(userID, accountID string, ignoreUnprocessedBefore time.Time) (int, error) {
	var count int
	// An unprocessed transaction is one that has at least one movement with an empty AccountId.
	// We also filter by accountID being present in at least one movement.
	var transactions []models.Transaction
	query := s.db.Where("user_id = ? AND movements LIKE ? AND merged_into_id IS NULL", userID, "%"+accountID+"%")
	if !ignoreUnprocessedBefore.IsZero() {
		query = query.Where("date >= ?", ignoreUnprocessedBefore)
	}
	err := query.Find(&transactions).Error
	if err != nil {
		return 0, fmt.Errorf(StorageError, err)
	}

	s.log.Debug("CountUnprocessedTransactionsForAccount query result",
		"accountId", accountID, "ignoreUnprocessedBefore", ignoreUnprocessedBefore, "totalTransactions", len(transactions))

	for _, t := range transactions {
		hasEmpty := false
		hasAccount := false
		for _, m := range t.Movements {
			// If a movement has 0 amount, it doesn't represent a financial impact
			// and shouldn't block reconciliation even if its AccountId is empty.
			if m.Amount.IsZero() {
				continue
			}
			if m.AccountId == "" {
				hasEmpty = true
			}
			if m.AccountId == accountID {
				hasAccount = true
			}
		}
		if hasEmpty && hasAccount {
			count++
		}
	}

	return count, nil
}

func (s *storage) HasTransactionsAfterDate(userID, accountID string, date time.Time) (bool, error) {
	var count int64
	err := s.db.Model(&models.Transaction{}).
		Where("user_id = ? AND movements LIKE ? AND date > ? AND merged_into_id IS NULL", userID, "%"+accountID+"%", date).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf(StorageError, err)
	}
	return count > 0, nil
}

func (s *storage) invalidateReconciliationIfAmountsChanged(
	userID string,
	oldMovements, newMovements []goserver.Movement,
	txDate time.Time,
) {
	// Build lookup for old movements
	oldByKey := make(map[string]goserver.Movement)
	for _, m := range oldMovements {
		key := m.AccountId + "|" + m.CurrencyId
		oldByKey[key] = m
	}

	// Check if any financial data changed
	affectedAccounts := make(map[string]string) // accountId -> currencyId

	for _, newM := range newMovements {
		if newM.AccountId == "" {
			continue // Unprocessed movements don't affect reconciliation directly
		}
		key := newM.AccountId + "|" + newM.CurrencyId
		oldM, exists := oldByKey[key]

		// New movement or amount changed
		if !exists || oldM.Amount != newM.Amount {
			affectedAccounts[newM.AccountId] = newM.CurrencyId
		}
		delete(oldByKey, key)
	}

	// Remaining old movements were removed
	for _, oldM := range oldByKey {
		if oldM.AccountId != "" {
			affectedAccounts[oldM.AccountId] = oldM.CurrencyId
		}
	}

	// Only invalidate if there were actual financial changes
	for accountId, currencyId := range affectedAccounts {
		lastRec, err := s.GetLatestReconciliation(userID, accountId, currencyId)
		if err != nil || lastRec == nil {
			continue
		}
		if txDate.Before(lastRec.ReconciledAt) {
			s.log.Info("Invalidating reconciliation due to financial change",
				"accountId", accountId, "currencyId", currencyId, "txDate", txDate, "recAt", lastRec.ReconciledAt)

			if err := s.InvalidateReconciliation(userID, accountId, currencyId, txDate); err != nil {
				s.log.Error("Failed to invalidate reconciliation", "error", err)
				continue
			}

			accountName := accountId
			if acc, err := s.GetAccount(userID, accountId); err == nil {
				accountName = acc.Name
			}

			_, _ = s.CreateNotification(userID, &goserver.Notification{
				Date:  time.Now(),
				Type:  string(models.NotificationTypeInfo),
				Title: "Reconciliation Invalidated",
				Description: fmt.Sprintf("Financial change to transaction before checkpoint invalidated reconciliation for account %q",
					accountName),
			})
		}
	}
}
