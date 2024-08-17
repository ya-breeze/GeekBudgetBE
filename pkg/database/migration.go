package database

import (
	"gorm.io/gorm"

	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
)

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Currency{},
		&models.Transaction{},
		&models.Matcher{},
		&models.BankImporter{},
	)
}
