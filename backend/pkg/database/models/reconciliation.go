package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

// Reconciliation tracks reconciliation checkpoints for accounts
type Reconciliation struct {
	gorm.Model

	UserID     string    `gorm:"index"`
	ID         uuid.UUID `gorm:"type:uuid;primaryKey"`
	AccountID  uuid.UUID `gorm:"type:uuid;index"`
	CurrencyID string    `gorm:"index"`

	ReconciledBalance decimal.Decimal
	ReconciledAt      time.Time
	ExpectedBalance   decimal.Decimal // Balance from bank importer or manually set
	IsManual          bool            // True for manual reconciliation (no bank importer)
}

// ReconciliationToAPI converts database model to API component
func ReconciliationToAPI(m *Reconciliation) *goserver.Reconciliation {
	return &goserver.Reconciliation{
		ReconciliationId:  m.ID.String(),
		AccountId:         m.AccountID.String(),
		CurrencyId:        m.CurrencyID,
		ReconciledBalance: m.ReconciledBalance,
		ReconciledAt:      m.ReconciledAt,
		ExpectedBalance:   m.ExpectedBalance,
		IsManual:          m.IsManual,
	}
}

// ReconciliationFromAPI converts API component to database model
func ReconciliationFromAPI(userID string, rec *goserver.ReconciliationNoId) *Reconciliation {
	accountId, _ := uuid.Parse(rec.AccountId)
	return &Reconciliation{
		UserID:            userID,
		AccountID:         accountId,
		CurrencyID:        rec.CurrencyId,
		ReconciledBalance: rec.ReconciledBalance,
		ExpectedBalance:   rec.ExpectedBalance,
		IsManual:          rec.IsManual,
	}
}
