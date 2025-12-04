package models

import (
	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type Currency struct {
	gorm.Model

	goserver.CurrencyNoId

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (a *Currency) FromDB() goserver.Currency {
	return goserver.Currency{
		Id:          a.ID.String(),
		Name:        a.Name,
		Description: a.Description,
	}
}

func (a *Currency) WithoutID() *goserver.CurrencyNoId {
	return &goserver.CurrencyNoId{
		Name:        a.Name,
		Description: a.Description,
	}
}

func CurrencyWithoutID(currency *goserver.Currency) *goserver.CurrencyNoId {
	return &goserver.CurrencyNoId{
		Name:        currency.Name,
		Description: currency.Description,
	}
}
