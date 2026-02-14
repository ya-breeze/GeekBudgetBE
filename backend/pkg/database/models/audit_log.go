package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuditLog struct {
	gorm.Model

	ID           uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID       string    `gorm:"index"`
	EntityType   string    `gorm:"index"` // e.g., "Transaction", "Account", "Matcher"
	EntityID     string    `gorm:"index"`
	Action       string    // e.g., "CREATED", "UPDATED", "DELETED", "MERGED"
	ChangeSource string    // e.g., "user", "system"

	// Before state of the entity (JSON).
	// For creates, this is null/empty.
	// For updates/deletes, this is the state before the action.
	Before *string `gorm:"type:text"`

	// After state of the entity (JSON).
	// For deletes, this is null/empty.
	// For creates/updates, this is the state after the action.
	After *string `gorm:"type:text"`

	CreatedAt time.Time
}
