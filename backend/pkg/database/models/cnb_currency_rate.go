package models

import (
	"time"
)

// CNBCurrencyRate represents currency rate for a specific date from Czech National Bank
type CNBCurrencyRate struct {
	CurrencyCode string    `gorm:"not null"`
	RateToCZK    float64   `gorm:"not null"`
	RateDate     time.Time `gorm:"not null;index"`
}
