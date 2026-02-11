package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type BudgetItem struct {
	gorm.Model

	Date        time.Time
	AccountID   string
	Amount      decimal.Decimal `gorm:"type:decimal(20,8)"`
	Description string

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (a *BudgetItem) FromDB() goserver.BudgetItem {
	return goserver.BudgetItem{
		Id:          a.ID.String(),
		Date:        a.Date,
		AccountId:   a.AccountID,
		Amount:      a.Amount,
		Description: a.Description,
	}
}

func BudgetItemToDB(m goserver.BudgetItemNoIdInterface, userID string) *BudgetItem {
	return &BudgetItem{
		UserID:      userID,
		Date:        m.GetDate(),
		AccountID:   m.GetAccountId(),
		Amount:      m.GetAmount(),
		Description: m.GetDescription(),
	}
}

func BudgetItemAccountWithoutID(budgetItem *goserver.BudgetItem) *goserver.BudgetItemNoId {
	return &goserver.BudgetItemNoId{
		Date:        budgetItem.Date,
		AccountId:   budgetItem.AccountId,
		Amount:      budgetItem.Amount,
		Description: budgetItem.Description,
	}
}
