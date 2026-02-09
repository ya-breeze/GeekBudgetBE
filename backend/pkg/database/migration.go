package database

import (
	"gorm.io/gorm"

	"github.com/ya-breeze/geekbudgetbe/pkg/database/models"
)

func autoMigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Account{},
		&models.Currency{},
		&models.Transaction{},
		&models.TransactionHistory{},

		&models.Matcher{},
		&models.BankImporter{},
		&models.Notification{},
		&models.Image{},
		&models.CNBCurrencyRate{},
		&models.BudgetItem{},
		&models.Reconciliation{},
		&models.TransactionDuplicate{},
	)
}
