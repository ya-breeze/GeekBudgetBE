package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type BankImporter struct {
	gorm.Model

	Name                 string
	Description          string
	AccountID            string
	Extra                string
	Type                 string
	LastSuccessfulImport time.Time
	LastImports          []goserver.BankImporterNoIdLastImportsInner `gorm:"serializer:json"`

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (t *BankImporter) FromDB() goserver.BankImporter {
	return goserver.BankImporter{
		Id:                   t.ID.String(),
		Name:                 t.Name,
		Description:          t.Description,
		AccountId:            t.AccountID,
		Extra:                t.Extra,
		Type:                 t.Type,
		LastSuccessfulImport: t.LastSuccessfulImport,
		LastImports:          t.LastImports,
	}
}

func BankImporterToDB(m *goserver.BankImporterNoId, userID string) *BankImporter {
	return &BankImporter{
		UserID:               userID,
		Name:                 m.Name,
		Description:          m.Description,
		AccountID:            m.AccountId,
		Extra:                m.Extra,
		Type:                 m.Type,
		LastSuccessfulImport: m.LastSuccessfulImport,
		LastImports:          m.LastImports,
	}
}
