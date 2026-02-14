package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

func (s *storage) GetAccounts(userID string) ([]goserver.Account, error) {
	result, err := s.db.Model(&models.Account{}).Where("user_id = ?", userID).Order("type, name").Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	accounts := make([]goserver.Account, 0)
	for result.Next() {
		var acc models.Account
		if err := s.db.ScanRows(result, &acc); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		accounts = append(accounts, acc.FromDB())
	}

	return accounts, nil
}

func (s *storage) CreateAccount(userID string, account *goserver.AccountNoId) (goserver.Account, error) {
	acc := models.AccountToDB(account, userID)
	acc.ID = uuid.New()
	if err := s.db.Create(&acc).Error; err != nil {
		return goserver.Account{}, fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, userID, "Account", acc.ID.String(), "CREATED", &acc); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	return acc.FromDB(), nil
}

func (s *storage) UpdateAccount(userID string, id string, account *goserver.AccountNoId) (goserver.Account, error) {
	return performUpdate[models.Account, goserver.AccountNoIdInterface, goserver.Account](s, userID, "Account", id, account,
		models.AccountToDB,
		func(m *models.Account) goserver.Account { return m.FromDB() },
		func(m *models.Account, id uuid.UUID) { m.ID = id },
	)
}

func (s *storage) DeleteAccount(userID string, id string, replaceWithAccountID *string) error {
	var acc models.Account
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&acc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrNotFound
		}
		return fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, userID, "Account", id, "DELETED", &acc); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if replaceWithAccountID != nil && *replaceWithAccountID != "" {
			newAccountID := *replaceWithAccountID

			// 1. Reassign BankImporters
			if err := tx.Model(&models.BankImporter{}).Where("account_id = ? AND user_id = ?", id, userID).
				Update("account_id", newAccountID).Error; err != nil {
				return fmt.Errorf("failed to reassign bank importers: %w", err)
			}

			// 2. Reassign Matchers
			if err := tx.Model(&models.Matcher{}).Where("output_account_id = ? AND user_id = ?", id, userID).
				Update("output_account_id", newAccountID).Error; err != nil {
				return fmt.Errorf("failed to reassign matchers: %w", err)
			}

			// 3. Reassign BudgetItems
			if err := tx.Model(&models.BudgetItem{}).Where("account_id = ? AND user_id = ?", id, userID).
				Update("account_id", newAccountID).Error; err != nil {
				return fmt.Errorf("failed to reassign budget items: %w", err)
			}

			// 4. Reassign Transactions (Movements)
			// Since movements are stored as JSON, we need to fetch, modify, and save.
			// Ideally we should process in batches, but for simplicity/MVP we do it here.
			// We find all transactions that have a movement with this account ID.
			// Optimization: Use SQLite's JSON functions for accurate querying.
			var transactions []models.Transaction
			if err := tx.Joins("CROSS JOIN json_each(transactions.movements)").
				Where("transactions.user_id = ? AND json_extract(json_each.value, '$.accountId') = ?", userID, id).
				Group("transactions.id").
				Find(&transactions).Error; err != nil {
				return fmt.Errorf("failed to find transactions for reassignment: %w", err)
			}

			for _, t := range transactions {
				updated := false
				newMovements := make([]goserver.Movement, len(t.Movements))
				for i, m := range t.Movements {
					if m.AccountId == id {
						m.AccountId = newAccountID
						updated = true
					}
					newMovements[i] = m
				}

				if updated {
					t.Movements = newMovements
					if err := tx.Save(&t).Error; err != nil {
						return fmt.Errorf("failed to save reassigned transaction %s: %w", t.ID, err)
					}
				}
			}
		} else {
			// User chose NOT to reassign.
			// Check if account is in use by any entity
			var count int64

			// Check BankImporters
			if err := tx.Model(&models.BankImporter{}).Where("account_id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check bank importers: %w", err)
			}
			if count > 0 {
				return ErrAccountInUse
			}

			// Check Matchers
			if err := tx.Model(&models.Matcher{}).Where("output_account_id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check matchers: %w", err)
			}
			if count > 0 {
				return ErrAccountInUse
			}

			// Check BudgetItems
			if err := tx.Model(&models.BudgetItem{}).Where("account_id = ? AND user_id = ?", id, userID).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check budget items: %w", err)
			}
			if count > 0 {
				return ErrAccountInUse
			}

			// Check Transactions (Movements)
			// Using SQLite's JSON functions for accurate checking
			if err := tx.Table("transactions").
				Joins("CROSS JOIN json_each(transactions.movements)").
				Where("transactions.user_id = ? AND json_extract(json_each.value, '$.accountId') = ?", userID, id).
				Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check transactions: %w", err)
			}
			if count > 0 {
				return ErrAccountInUse
			}
		}

		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Account{}).Error; err != nil {
			return fmt.Errorf(StorageError, err)
		}
		return nil
	})
}

func (s *storage) GetAccountHistory(userID string, accountID string) ([]goserver.Transaction, error) {
	// result, err := s.db.Model(&models.Transaction{}).Where("user_id = ? AND account_id = ?", userID, accountID).Rows()
	// if err != nil {
	// 	return nil, fmt.Errorf(StorageError, err)
	// }
	// defer result.Close()

	var transactions []goserver.Transaction
	// for result.Next() {
	// 	var tr models.Transaction
	// 	if err := s.db.ScanRows(result, &tr); err != nil {
	// 		return nil, fmt.Errorf(StorageError, err)
	// 	}
	//
	// 	transactions = append(transactions, tr.FromDB())
	// }

	return transactions, nil
}

func (s *storage) GetAccount(userID string, id string) (goserver.Account, error) {
	var acc models.Account
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&acc).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Account{}, ErrNotFound
		}

		return goserver.Account{}, fmt.Errorf(StorageError, err)
	}

	return acc.FromDB(), nil
}

func (s *storage) GetAccountBalance(userID, accountID, currencyID string) (decimal.Decimal, error) {
	acc, err := s.GetAccount(userID, accountID)
	if err != nil {
		return decimal.Zero, err
	}

	var total decimal.Decimal
	for _, b := range acc.BankInfo.Balances {
		if b.CurrencyId == currencyID {
			total = total.Add(b.OpeningBalance)
			break
		}
	}

	// Sum all movements for this account and currency
	// We use raw SQL to iterate over the movements JSON column in SQLite
	// Since movements is a JSON array of objects, we need to parse it.
	// For simplicity and to avoid complex SQLite JSON path expressions that might vary,
	// we fetch transactions and sum in Go, but filter at DB level if possible.
	var transactions []models.Transaction
	err = s.db.Where("user_id = ? AND movements LIKE ? AND merged_into_id IS NULL", userID, "%"+accountID+"%").Find(&transactions).Error
	if err != nil {
		return decimal.Zero, fmt.Errorf(StorageError, err)
	}

	for _, t := range transactions {
		for _, m := range t.Movements {
			if m.AccountId == accountID && m.CurrencyId == currencyID {
				total = total.Add(m.Amount)
			}
		}
	}

	return total, nil
}
