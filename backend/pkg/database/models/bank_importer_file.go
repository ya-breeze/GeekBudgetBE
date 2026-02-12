package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type BankImporterFile struct {
	gorm.Model

	UserID         string    `gorm:"index"`
	BankImporterID uuid.UUID `gorm:"type:uuid;index"`
	Filename       string
	Path           string // Relative path to the stored file
	UploadDate     time.Time
	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (t *BankImporterFile) FromDB() goserver.BankImporterFile {
	return goserver.BankImporterFile{
		Id:             t.ID.String(),
		BankImporterId: t.BankImporterID.String(),
		Filename:       t.Filename,
		UploadDate:     t.UploadDate,
	}
}
