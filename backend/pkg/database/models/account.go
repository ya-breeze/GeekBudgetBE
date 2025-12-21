package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type Account struct {
	gorm.Model

	Name                   string
	Description            string
	Type                   string
	ShowInDashboardSummary bool
	HideFromReports        bool `gorm:"default:false"`
	Image                  string

	BankInfo goserver.BankAccountInfo `gorm:"serializer:json"`

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (a *Account) FromDB() goserver.Account {
	return goserver.Account{
		Id:                     a.ID.String(),
		Name:                   a.Name,
		Type:                   a.Type,
		Description:            a.Description,
		BankInfo:               a.BankInfo,
		ShowInDashboardSummary: a.ShowInDashboardSummary,
		HideFromReports:        a.HideFromReports,
		Image:                  a.Image,
	}
}

func AccountToDB(m goserver.AccountNoIdInterface, userID string) *Account {
	return &Account{
		UserID:                 userID,
		Name:                   m.GetName(),
		Description:            m.GetDescription(),
		Type:                   m.GetType(),
		BankInfo:               m.GetBankInfo(),
		ShowInDashboardSummary: m.GetShowInDashboardSummary(),
		HideFromReports:        m.GetHideFromReports(),
		Image:                  m.GetImage(),
	}
}

func AccountWithoutID(account *goserver.Account) *goserver.AccountNoId {
	return &goserver.AccountNoId{
		Name:                   account.Name,
		Type:                   account.Type,
		Description:            account.Description,
		BankInfo:               account.BankInfo,
		ShowInDashboardSummary: account.ShowInDashboardSummary,
		HideFromReports:        account.HideFromReports,
		Image:                  account.Image,
	}
}
