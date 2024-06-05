package database

import (
	"gorm.io/gorm"

	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
)

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Account{},
	)
}
