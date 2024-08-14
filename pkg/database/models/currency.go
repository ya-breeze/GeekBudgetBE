package models

import (
	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type Currency struct {
	gorm.Model

	goserver.CurrencyNoId

	UserID string    `gorm:"index:idx_user_id"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (a *Currency) FromDB() goserver.Currency {
	return goserver.Currency{
		Id:          a.ID.String(),
		Name:        a.Name,
		Description: a.Description,
	}
}
