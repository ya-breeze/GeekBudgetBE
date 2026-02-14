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

	// Snapshot of the entity at the time of the action.
	// For updates, this is usually the new state.
	// For deletes, this is the state before deletion.
	Snapshot string `gorm:"type:text"`

	CreatedAt time.Time
}
