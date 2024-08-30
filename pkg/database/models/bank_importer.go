package models

import (
	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type BankImporter struct {
	gorm.Model

	Name        string
	Description string
	Extra       string

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (t *BankImporter) FromDB() goserver.BankImporter {
	return goserver.BankImporter{
		Id:          t.ID.String(),
		Name:        t.Name,
		Description: t.Description,
		Extra:       t.Extra,
	}
}

func BankImporterToDB(m *goserver.BankImporterNoId, userID string) *BankImporter {
	return &BankImporter{
		UserID:      userID,
		Name:        m.Name,
		Description: m.Description,
		Extra:       m.Extra,
	}
}
