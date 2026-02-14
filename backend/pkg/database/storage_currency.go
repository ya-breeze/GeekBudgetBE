package database

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

func (s *storage) CreateCurrency(userID string, currency *goserver.CurrencyNoId) (goserver.Currency, error) {
	cur := models.Currency{
		ID:           uuid.New(),
		UserID:       userID,
		CurrencyNoId: *currency,
	}
	if err := s.db.Create(&cur).Error; err != nil {
		return goserver.Currency{}, fmt.Errorf(StorageError, err)
	}

	if err := s.recordAuditLog(s.db, userID, "Currency", cur.ID.String(), "CREATED", &cur); err != nil {
		s.log.Error("Failed to record audit log", "error", err)
	}

	return cur.FromDB(), nil
}

func (s *storage) GetCurrencies(userID string) ([]goserver.Currency, error) {
	result, err := s.db.Model(&models.Currency{}).Where("user_id = ?", userID).Rows()
	if err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}
	defer result.Close()

	currencies := make([]goserver.Currency, 0)
	for result.Next() {
		var cur models.Currency
		if err := s.db.ScanRows(result, &cur); err != nil {
			return nil, fmt.Errorf(StorageError, err)
		}

		currencies = append(currencies, cur.FromDB())
	}

	return currencies, nil
}

func (s *storage) GetCurrency(userID string, id string) (goserver.Currency, error) {
	var cur models.Currency
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&cur).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return goserver.Currency{}, ErrNotFound
		}

		return goserver.Currency{}, fmt.Errorf(StorageError, err)
	}

	return cur.FromDB(), nil
}

func (s *storage) UpdateCurrency(userID string, id string, currency *goserver.CurrencyNoId) (goserver.Currency, error) {
	return performUpdate[models.Currency, *goserver.CurrencyNoId, goserver.Currency](s, userID, "Currency", id, currency,
		func(i *goserver.CurrencyNoId, userID string) *models.Currency {
			return &models.Currency{UserID: userID, CurrencyNoId: *i}
		},
		func(m *models.Currency) goserver.Currency { return m.FromDB() },
		func(m *models.Currency, id uuid.UUID) { m.ID = id },
	)
}

func (s *storage) DeleteCurrency(userID string, id string, replaceWithCurrencyID *string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		if replaceWithCurrencyID != nil && *replaceWithCurrencyID != "" {
			newCurrencyID := *replaceWithCurrencyID

			// 1. Reassign in User (favorite currency)
			if err := tx.Model(&models.User{}).Where("id = ? AND favorite_currency_id = ?", userID, id).
				Update("favorite_currency_id", newCurrencyID).Error; err != nil {
				return fmt.Errorf("failed to reassign user favorite currency: %w", err)
			}

			// 2. Reassign in Accounts (BankInfo)
			var accounts []models.Account
			if err := tx.Joins("CROSS JOIN json_each(accounts.bank_info, '$.balances')").
				Where("accounts.user_id = ? AND json_extract(json_each.value, '$.currencyId') = ?", userID, id).
				Group("accounts.id").
				Find(&accounts).Error; err != nil {
				return fmt.Errorf("failed to find accounts for currency reassignment: %w", err)
			}
			for _, acc := range accounts {
				updated := false
				newBalances := make([]goserver.BankAccountInfoBalancesInner, len(acc.BankInfo.Balances))

				for i, b := range acc.BankInfo.Balances {
					if b.CurrencyId == id {
						b.CurrencyId = newCurrencyID
						updated = true
					}
					newBalances[i] = b
				}

				if updated {
					acc.BankInfo.Balances = newBalances
					if err := tx.Save(&acc).Error; err != nil {
						return fmt.Errorf("failed to save reassigned account %s: %w", acc.ID, err)
					}
				}
			}

			// 3. Reassign in Transactions (Movements)
			var transactions []models.Transaction
			if err := tx.Joins("CROSS JOIN json_each(transactions.movements)").
				Where("transactions.user_id = ? AND json_extract(json_each.value, '$.currencyId') = ?", userID, id).
				Group("transactions.id").
				Find(&transactions).Error; err != nil {
				return fmt.Errorf("failed to find transactions for currency reassignment: %w", err)
			}

			for _, t := range transactions {
				updated := false
				newMovements := make([]goserver.Movement, len(t.Movements))
				for i, m := range t.Movements {
					if m.CurrencyId == id {
						m.CurrencyId = newCurrencyID
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
			// Check if currency is in use
			var count int64

			// Check User favorite currency
			if err := tx.Model(&models.User{}).Where("id = ? AND favorite_currency_id = ?", userID, id).Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check user favorite currency: %w", err)
			}
			if count > 0 {
				return ErrCurrencyInUse
			}

			// Check Accounts
			// Using SQLite's JSON functions for accurate checking
			if err := tx.Table("accounts").
				Joins("CROSS JOIN json_each(accounts.bank_info, '$.balances')").
				Where("accounts.user_id = ? AND json_extract(json_each.value, '$.currencyId') = ?", userID, id).
				Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check accounts for currency usage: %w", err)
			}
			if count > 0 {
				return ErrCurrencyInUse
			}

			// Check Transactions
			// Using SQLite's JSON functions for accurate checking
			if err := tx.Table("transactions").
				Joins("CROSS JOIN json_each(transactions.movements)").
				Where("transactions.user_id = ? AND json_extract(json_each.value, '$.currencyId') = ?", userID, id).
				Count(&count).Error; err != nil {
				return fmt.Errorf("failed to check transactions for currency usage: %w", err)
			}
			if count > 0 {
				return ErrCurrencyInUse
			}
		}

		var cur models.Currency
		if err := tx.Where("id = ? AND user_id = ?", id, userID).First(&cur).Error; err == nil {
			if err := s.recordAuditLog(tx, userID, "Currency", id, "DELETED", &cur); err != nil {
				s.log.Error("Failed to record audit log", "error", err)
			}
		}

		if err := tx.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Currency{}).Error; err != nil {
			return fmt.Errorf(StorageError, err)
		}
		return nil
	})
}
