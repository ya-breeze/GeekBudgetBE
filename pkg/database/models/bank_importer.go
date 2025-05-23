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
	LastImports          []goserver.ImportResult                  `gorm:"serializer:json"`
	Mappings             []goserver.BankImporterNoIdMappingsInner `gorm:"serializer:json"`
	FetchAll             bool
	UserID               string    `gorm:"index"`
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey"`
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
		Mappings:             t.Mappings,
		FetchAll:             t.FetchAll,
	}
}

func BankImporterToDB(m goserver.BankImporterNoIdInterface, userID string) *BankImporter {
	return &BankImporter{
		UserID:               userID,
		Name:                 m.GetName(),
		Description:          m.GetDescription(),
		AccountID:            m.GetAccountId(),
		Extra:                m.GetExtra(),
		Type:                 m.GetType(),
		LastSuccessfulImport: m.GetLastSuccessfulImport(),
		LastImports:          m.GetLastImports(),
		FetchAll:             m.GetFetchAll(),
	}
}

func BankImporterWithoutID(bankImporter *goserver.BankImporter) *goserver.BankImporterNoId {
	return &goserver.BankImporterNoId{
		Name:                 bankImporter.Name,
		Description:          bankImporter.Description,
		AccountId:            bankImporter.AccountId,
		Extra:                bankImporter.Extra,
		Type:                 bankImporter.Type,
		LastSuccessfulImport: bankImporter.LastSuccessfulImport,
		LastImports:          bankImporter.LastImports,
		FetchAll:             bankImporter.FetchAll,
	}
}
