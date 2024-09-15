package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type Account struct {
	gorm.Model

	Name        string
	Description string
	Type        string

	BankInfo goserver.BankAccountInfo `gorm:"serializer:json"`

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (a *Account) FromDB() goserver.Account {
	return goserver.Account{
		Id:          a.ID.String(),
		Name:        a.Name,
		Type:        a.Type,
		Description: a.Description,
		BankInfo:    a.BankInfo,
	}
}

func AccountToDB(m goserver.AccountNoIdInterface, userID string) *Account {
	return &Account{
		UserID:      userID,
		Name:        m.GetName(),
		Description: m.GetDescription(),
		Type:        m.GetType(),
		BankInfo:    m.GetBankInfo(),
	}
}

func AccountWithoutID(account *goserver.Account) *goserver.AccountNoId {
	return &goserver.AccountNoId{
		Name:        account.Name,
		Type:        account.Type,
		Description: account.Description,
		BankInfo:    account.BankInfo,
	}
}
