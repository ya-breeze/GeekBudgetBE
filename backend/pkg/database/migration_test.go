package database

import (
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
)

func TestMigration(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}

	// Open database manually to control migration
	db, err := openSqlite(logger, cfg.DBPath, cfg.Verbose)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// 1. Setup initial schema without MergedTransaction
	if err := db.AutoMigrate(
		&models.User{},
		&models.Transaction{},
	); err != nil {
		t.Fatalf("failed to migrate initial schema: %v", err)
	}

	userID := "user-1"
	kID := uuid.New()
	tID := uuid.New()
	now := time.Now()

	// 2. Insert transactions representing an old-style merge (soft-deleted, with merged_into_id)
	k := models.Transaction{
		ID:          kID,
		UserID:      userID,
		Date:        now,
		Description: "Keep",
	}
	m := models.Transaction{
		ID:           tID,
		UserID:       userID,
		Date:         now,
		Description:  "Merge",
		MergedIntoID: &kID,
		MergedAt:     &now,
	}

	if err := db.Create(&k).Error; err != nil {
		t.Fatalf("failed to create keep transaction: %v", err)
	}
	if err := db.Create(&m).Error; err != nil {
		t.Fatalf("failed to create merge transaction: %v", err)
	}
	// Soft delete it the old way
	if err := db.Delete(&m).Error; err != nil {
		t.Fatalf("failed to soft-delete: %v", err)
	}

	// 3. Insert another one that is merged but NOT soft-deleted (edge case from user feedback)
	tID2 := uuid.New()
	m2 := models.Transaction{
		ID:           tID2,
		UserID:       userID,
		Date:         now,
		Description:  "Merge 2",
		MergedIntoID: &kID,
		MergedAt:     &now,
	}
	if err := db.Create(&m2).Error; err != nil {
		t.Fatalf("failed to create merge 2 transaction: %v", err)
	}

	// 4. Run the full migration including archive creation
	if err := autoMigrateModels(db); err != nil {
		t.Fatalf("failed to run full migration: %v", err)
	}

	// 5. Verify archive table contains both
	var count int64
	if err := db.Model(&models.MergedTransaction{}).Count(&count).Error; err != nil {
		t.Fatalf("failed to count merged transactions: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 archived transactions, got %d", count)
	}

	// 6. Verify original transactions are completely removed (hard-deleted) from transactions table
	var t1Remaining, t2Remaining int64
	db.Unscoped().Model(&models.Transaction{}).Where("id = ?", tID).Count(&t1Remaining)
	db.Unscoped().Model(&models.Transaction{}).Where("id = ?", tID2).Count(&t2Remaining)

	if t1Remaining != 0 || t2Remaining != 0 {
		t.Errorf("expected original transactions to be hard-deleted (not exist even with Unscoped), got %d and %d", t1Remaining, t2Remaining)
	}
}
