package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite(dbPath string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
}
