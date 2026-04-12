package database

import (
	"encoding/base64"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	kinauth "github.com/ya-breeze/kin-core/auth"
	"golang.org/x/crypto/bcrypt"

	"github.com/ya-breeze/geekbudgetbe/pkg/config"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
)

// TestKinCoreMigrationPasswordCompat verifies that users migrated from the old schema
// (which stored base64(bcrypt) passwords) can still log in after migration.
// The old auth decoded base64 before calling bcrypt.CompareHashAndPassword.
// kin-core's VerifyPassword expects raw bcrypt strings, so the migration must
// strip the base64 encoding layer.
func TestKinCoreMigrationPasswordCompat(t *testing.T) {
	logger := slog.Default()
	cfg := &config.Config{DBPath: ":memory:", Verbose: false}

	db, err := openSqlite(logger, cfg.DBPath, cfg.Verbose)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}

	// Create the OLD schema — users table with 'login' column triggers migration detection.
	// All other user-scoped tables must exist (even empty) because the migration reads them.
	oldTables := []string{
		`CREATE TABLE users (id TEXT PRIMARY KEY, login TEXT UNIQUE, hashed_password TEXT,
			start_date DATETIME, favorite_currency_id TEXT,
			created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE accounts (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE currencies (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE transactions (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE transaction_histories (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE audit_logs (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE matchers (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE bank_importers (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE bank_importer_files (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE notifications (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE budget_items (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE reconciliations (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE transaction_duplicates (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE merged_transactions (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
		`CREATE TABLE transaction_templates (id TEXT PRIMARY KEY, user_id TEXT, created_at DATETIME, updated_at DATETIME, deleted_at DATETIME)`,
	}
	for _, ddl := range oldTables {
		if err := db.Exec(ddl).Error; err != nil {
			t.Fatalf("failed to create old table: %v\nDDL: %s", err, ddl)
		}
	}

	// Simulate old password storage: bcrypt the password, then base64-encode it
	plainPassword := "my-real-password"
	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("failed to bcrypt password: %v", err)
	}
	oldStoredHash := base64.StdEncoding.EncodeToString(bcryptHash) // old format: base64(bcrypt)

	userID := uuid.New().String()
	now := time.Now()
	if err := db.Exec(`INSERT INTO users (id, login, hashed_password, start_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`, userID, "testuser", oldStoredHash, now, now, now).Error; err != nil {
		t.Fatalf("failed to insert old user: %v", err)
	}

	// Run kin-core migration
	if err := runMigrationIfNeeded(logger, db); err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	// After migration, find the user and verify the password works with kin-core
	var migratedUser models.User
	if err := db.Where("username = ?", "testuser").First(&migratedUser).Error; err != nil {
		t.Fatalf("migrated user not found: %v", err)
	}

	// kin-core VerifyPassword must work with the migrated hash
	if !kinauth.VerifyPassword(plainPassword, migratedUser.PasswordHash) {
		t.Errorf("VerifyPassword failed after migration — existing users cannot log in")
	}

	// Wrong password must still fail
	if kinauth.VerifyPassword("wrong-password", migratedUser.PasswordHash) {
		t.Errorf("VerifyPassword accepted wrong password")
	}
}

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

	familyID := uuid.MustParse("00000000-0000-0000-0000-000000000001")
	kID := uuid.New()
	tID := uuid.New()
	now := time.Now()

	// 2. Insert transactions representing an old-style merge (soft-deleted, with merged_into_id)
	k := models.Transaction{
		ID:       kID,
		FamilyID: familyID,
		Date:     now,
		Description: "Keep",
	}
	m := models.Transaction{
		ID:           tID,
		FamilyID:     familyID,
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
		FamilyID:     familyID,
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
