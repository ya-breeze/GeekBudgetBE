package database

import (
	"database/sql"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
	"github.com/ya-breeze/kin-core/authdb"
	"gorm.io/gorm"
)

// oldUserRow represents the pre-kin-core user row for migration purposes.
type oldUserRow struct {
	ID                 string    `gorm:"column:id"`
	Login              string    `gorm:"column:login"`
	HashedPassword     string    `gorm:"column:hashed_password"`
	StartDate          time.Time `gorm:"column:start_date"`
	FavoriteCurrencyID string    `gorm:"column:favorite_currency_id"`
}

func (oldUserRow) TableName() string { return "users" }

// userIDToFamilyID maps old string user_id values to new family UUIDs.
type userIDToFamilyID map[string]uuid.UUID

// runMigrationIfNeeded detects old schema (login column in users) and migrates
// to family-based schema. Safe to call on every startup — returns immediately if
// the schema is already up to date.
func runMigrationIfNeeded(log *slog.Logger, db *gorm.DB) error {
	// Detect old schema: old users table had 'login' column; new schema uses 'username'.
	var count int64
	db.Raw("SELECT COUNT(*) FROM pragma_table_info('users') WHERE name='login'").Scan(&count)
	if count == 0 {
		return nil // already migrated or fresh DB
	}

	log.Info("Detected old schema (users.login column found) — starting kin-core migration")

	// Step 1: Read all old users into memory before dropping anything.
	var oldUsers []oldUserRow
	if err := db.Find(&oldUsers).Error; err != nil {
		return fmt.Errorf("read old users: %w", err)
	}

	// Build userID → familyID mapping (one family per user).
	mapping := make(userIDToFamilyID, len(oldUsers))
	for _, u := range oldUsers {
		mapping[u.ID] = uuid.New()
	}

	// Step 2: Read all 15 user-scoped tables into memory.
	tableData := make(map[string][][]interface{}) // table → rows (each row = ordered column values)
	tableColumns := make(map[string][]string)     // table → column names

	// Tables that have user_id and need family_id mapping.
	// (images and cnb_currency_rates are global — skipped)
	userScopedTables := []string{
		"accounts",
		"currencies",
		"transactions",
		"transaction_histories",
		"audit_logs",
		"matchers",
		"bank_importers",
		"bank_importer_files",
		"notifications",
		"budget_items",
		"reconciliations",
		"transaction_duplicates",
		"merged_transactions",
		"transaction_templates",
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}

	for _, table := range userScopedTables {
		cols, rows, err := readTableRows(sqlDB, table)
		if err != nil {
			return fmt.Errorf("read table %s: %w", table, err)
		}
		tableColumns[table] = cols
		tableData[table] = rows
		log.Info("Read table", "table", table, "rows", len(rows))
	}

	log.Info("Read old data", "users", len(oldUsers))

	// Step 3: Drop all old tables (including indexes — preventing name collision on recreate).
	allTables := append([]string{"families", "users"}, userScopedTables...)
	for _, table := range allTables {
		if err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %q", table)).Error; err != nil {
			return fmt.Errorf("drop table %s: %w", table, err)
		}
	}
	log.Info("Old tables dropped")

	// Step 4: Recreate schema via GORM AutoMigrate (exact match ensures no re-migration next time).
	if err := db.AutoMigrate(
		&models.Family{},
		&models.User{},
		&models.Account{},
		&models.Currency{},
		&models.Transaction{},
		&models.TransactionHistory{},
		&models.AuditLog{},
		&models.Matcher{},
		&models.BankImporter{},
		&models.Notification{},
		&models.BudgetItem{},
		&models.BankImporterFile{},
		&models.Reconciliation{},
		&models.TransactionDuplicate{},
		&models.MergedTransaction{},
		&models.TransactionTemplate{},
		&authdb.RefreshToken{},
		&authdb.BlacklistedToken{},
	); err != nil {
		return fmt.Errorf("auto-migrate new schema: %w", err)
	}
	log.Info("New schema created")

	// Step 5: Insert migrated data.
	now := time.Now()

	// Families: one per old user.
	for _, u := range oldUsers {
		familyID := mapping[u.ID]
		family := models.Family{}
		family.ID = familyID
		family.Name = u.Login
		family.CreatedAt = now
		family.UpdatedAt = now
		if err := db.Create(&family).Error; err != nil {
			return fmt.Errorf("insert family for %s: %w", u.Login, err)
		}
	}
	log.Info("Families inserted", "count", len(oldUsers))

	// Users: remap login→username, hashed_password→password_hash, add family_id.
	for _, u := range oldUsers {
		familyID := mapping[u.ID]
		user := models.User{}
		user.ID = uuid.MustParse(u.ID)
		user.Username = u.Login
		user.PasswordHash = u.HashedPassword
		user.FamilyID = familyID
		user.StartDate = u.StartDate
		user.FavoriteCurrencyID = u.FavoriteCurrencyID
		user.CreatedAt = now
		user.UpdatedAt = now
		if err := db.Create(&user).Error; err != nil {
			return fmt.Errorf("insert migrated user %s: %w", u.Login, err)
		}
	}
	log.Info("Users migrated", "count", len(oldUsers))

	// Data tables: copy all rows, replacing user_id with mapped family_id.
	for _, table := range userScopedTables {
		cols := tableColumns[table]
		rows := tableData[table]

		newCols := remapColumns(cols)

		inserted := 0
		for _, row := range rows {
			newRow := remapRow(cols, newCols, row, mapping)
			if err := insertRow(db, table, newCols, newRow); err != nil {
				return fmt.Errorf("insert row in %s: %w", table, err)
			}
			inserted++
		}
		log.Info("Table migrated", "table", table, "rows", inserted)
	}

	log.Info("Migration complete")
	return nil
}

// readTableRows reads all rows from a table. Returns column names and row values.
func readTableRows(sqlDB *sql.DB, table string) ([]string, [][]interface{}, error) {
	rows, err := sqlDB.Query(fmt.Sprintf("SELECT * FROM %q", table)) //nolint:gosec // table name is hardcoded
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, nil, err
	}

	var result [][]interface{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, nil, err
		}
		row := make([]interface{}, len(cols))
		copy(row, vals)
		result = append(result, row)
	}
	return cols, result, rows.Err()
}

// remapColumns returns a new column list with user_id → family_id.
func remapColumns(cols []string) []string {
	newCols := make([]string, len(cols))
	for i, c := range cols {
		if c == "user_id" {
			newCols[i] = "family_id"
		} else {
			newCols[i] = c
		}
	}
	return newCols
}

// remapRow returns a new row with the user_id value replaced by the mapped family_id UUID string.
func remapRow(oldCols, newCols []string, row []interface{}, mapping userIDToFamilyID) []interface{} {
	newRow := make([]interface{}, len(row))
	copy(newRow, row)
	for i, col := range oldCols {
		if col == "user_id" {
			if row[i] != nil {
				userID := fmt.Sprintf("%v", row[i])
				if familyID, ok := mapping[userID]; ok {
					newRow[i] = familyID.String()
				} else {
					// Unknown user_id — use zero UUID as fallback (shouldn't happen)
					newRow[i] = uuid.Nil.String()
				}
			}
		}
	}
	return newRow
}

// insertRow inserts a single row using raw SQL.
func insertRow(db *gorm.DB, table string, cols []string, vals []interface{}) error {
	quoted := make([]string, len(cols))
	for i, c := range cols {
		quoted[i] = fmt.Sprintf("%q", c)
	}
	placeholders := strings.Repeat("?,", len(cols))
	placeholders = placeholders[:len(placeholders)-1]

	query := fmt.Sprintf("INSERT INTO %q (%s) VALUES (%s)", table, strings.Join(quoted, ","), placeholders) //nolint:gosec // table name is hardcoded

	return db.Exec(query, vals...).Error
}
