package database

import (
	"errors"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"gorm.io/gorm"
)

func (s *storage) SaveCNBRates(rates map[string]decimal.Decimal, date time.Time) error {
	// Use a transaction to ensure all rates are saved together
	return s.db.Transaction(func(tx *gorm.DB) error {
		// First delete all existing rates for this date to avoid duplicates
		if err := tx.Where("rate_date = ?", date).Delete(&models.CNBCurrencyRate{}).Error; err != nil {
			return fmt.Errorf(StorageError, err)
		}

		// Create new rates
		for currencyCode, rate := range rates {
			currencyRate := models.CNBCurrencyRate{
				CurrencyCode: currencyCode,
				RateToCZK:    rate,
				RateDate:     date,
			}

			if err := tx.Create(&currencyRate).Error; err != nil {
				return fmt.Errorf(StorageError, err)
			}
		}

		return nil
	})
}

func (s *storage) GetCNBRates(date time.Time) (map[string]decimal.Decimal, error) {
	var rates []models.CNBCurrencyRate
	query := s.db.Model(&models.CNBCurrencyRate{})

	// If a specific date is provided, use it
	if !date.IsZero() {
		query = query.Where("rate_date = ?", date)
	} else {
		// Otherwise get the most recent rates
		var latestRate models.CNBCurrencyRate
		if err := s.db.Model(&models.CNBCurrencyRate{}).
			Order("rate_date DESC").
			First(&latestRate).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return make(map[string]decimal.Decimal), nil
			}
			return nil, fmt.Errorf(StorageError, err)
		}

		query = query.Where("rate_date = ?", latestRate.RateDate)
	}

	if err := query.Find(&rates).Error; err != nil {
		return nil, fmt.Errorf(StorageError, err)
	}

	// Convert to map
	result := make(map[string]decimal.Decimal, len(rates))
	for _, rate := range rates {
		result[rate.CurrencyCode] = rate.RateToCZK
	}

	return result, nil
}
