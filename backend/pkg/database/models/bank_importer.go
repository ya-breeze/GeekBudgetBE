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
	FeeAccountID         string
	LastSuccessfulImport time.Time
	LastImports          []goserver.ImportResult                  `gorm:"serializer:json"`
	Mappings             []goserver.BankImporterNoIdMappingsInner `gorm:"serializer:json"`
	FetchAll             bool
	IsStopped            bool
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
		FeeAccountId:         t.FeeAccountID,
		Type:                 t.Type,
		LastSuccessfulImport: t.LastSuccessfulImport,
		LastImports:          t.LastImports,
		Mappings:             t.Mappings,
		FetchAll:             t.FetchAll,
		IsStopped:            t.IsStopped,
	}
}

func BankImporterToDB(m goserver.BankImporterNoIdInterface, userID string) *BankImporter {
	return &BankImporter{
		UserID:               userID,
		Name:                 m.GetName(),
		Description:          m.GetDescription(),
		AccountID:            m.GetAccountId(),
		FeeAccountID:         m.GetFeeAccountId(),
		Extra:                m.GetExtra(),
		Type:                 m.GetType(),
		LastSuccessfulImport: m.GetLastSuccessfulImport(),
		LastImports:          m.GetLastImports(),
		FetchAll:             m.GetFetchAll(),
		IsStopped:            m.GetIsStopped(),
	}
}

func BankImporterWithoutID(bankImporter *goserver.BankImporter) *goserver.BankImporterNoId {
	return &goserver.BankImporterNoId{
		Name:                 bankImporter.Name,
		Description:          bankImporter.Description,
		AccountId:            bankImporter.AccountId,
		FeeAccountId:         bankImporter.FeeAccountId,
		Extra:                bankImporter.Extra,
		Type:                 bankImporter.Type,
		LastSuccessfulImport: bankImporter.LastSuccessfulImport,
		LastImports:          bankImporter.LastImports,
		FetchAll:             bankImporter.FetchAll,
		IsStopped:            bankImporter.IsStopped,
	}
}
