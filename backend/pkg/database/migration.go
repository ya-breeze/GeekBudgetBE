package database

import (
	"time"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
)

func autoMigrateModels(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Currency{},
		&models.Transaction{},
		&models.TransactionHistory{},
		&models.AuditLog{},

		&models.Matcher{},
		&models.BankImporter{},
		&models.Notification{},
		&models.Image{},
		&models.CNBCurrencyRate{},
		&models.BudgetItem{},
		&models.BankImporterFile{},
		&models.Reconciliation{},
		&models.TransactionDuplicate{},

		&models.MergedTransaction{},
	); err != nil {
		return err
	}

	return migrateExistingMergedTransactions(db)
}

func migrateExistingMergedTransactions(db *gorm.DB) error {
	var transactions []models.Transaction
	// Find all merged transactions (even if not yet soft-deleted)
	if err := db.Unscoped().Where("merged_into_id IS NOT NULL").Find(&transactions).Error; err != nil {
		return err
	}

	for _, t := range transactions {
		// Check if already archived
		var count int64
		if err := db.Model(&models.MergedTransaction{}).Where("original_transaction_id = ?", t.ID).Count(&count).Error; err != nil {
			return err
		}
		if count > 0 {
			continue
		}

		mergedAt := time.Now()
		if t.MergedAt != nil {
			mergedAt = *t.MergedAt
		}

		archive := models.MergedTransaction{
			ID:                    uuid.New(),
			UserID:                t.UserID,
			KeptTransactionID:     *t.MergedIntoID,
			OriginalTransactionID: t.ID,
			Date:                  t.Date,
			Description:           t.Description,
			Place:                 t.Place,
			Tags:                  t.Tags,
			PartnerName:           t.PartnerName,
			PartnerAccount:        t.PartnerAccount,
			PartnerInternalID:     t.PartnerInternalID,
			Extra:                 t.Extra,
			UnprocessedSources:    t.UnprocessedSources,
			ExternalIDs:           t.ExternalIDs,
			Movements:             t.Movements,
			MatcherID:             t.MatcherID,
			IsAuto:                t.IsAuto,
			SuspiciousReasons:     t.SuspiciousReasons,
			MergedAt:              mergedAt,
		}

		if err := db.Create(&archive).Error; err != nil {
			return err
		}
	}

	// Hard-delete all migrated transactions (archive is source of truth)
	return db.Unscoped().Where("merged_into_id IS NOT NULL").Delete(&models.Transaction{}).Error
}
