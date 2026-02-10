package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type MergedTransaction struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID                string    `gorm:"index"`
	KeptTransactionID     uuid.UUID `gorm:"type:uuid;index"`
	OriginalTransactionID uuid.UUID `gorm:"type:uuid;index"`

	// Snapshot of original transaction fields
	Date               time.Time
	Description        string
	Place              string
	Tags               []string `gorm:"serializer:json"`
	PartnerName        string
	PartnerAccount     string
	PartnerInternalID  string
	Extra              string
	UnprocessedSources string
	ExternalIDs        []string            `gorm:"serializer:json"`
	Movements          []goserver.Movement `gorm:"serializer:json"`
	MatcherID          *uuid.UUID          `gorm:"type:uuid"`
	IsAuto             bool
	SuspiciousReasons  []string `gorm:"serializer:json"`

	MergedAt  time.Time `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
