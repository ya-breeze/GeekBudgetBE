package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionHistory struct {
	gorm.Model

	TransactionID uuid.UUID `gorm:"type:uuid;index"`
	UserID        string    `gorm:"index"`
	Action        string
	// Snapshot of the transaction at the time of the action (or before for updates/deletes)
	// Stored as JSON to be resilient to future schema changes in Transaction model
	Snapshot string `gorm:"type:text"`

	CreatedAt time.Time
	ID        uuid.UUID `gorm:"type:uuid;primaryKey"`
}
