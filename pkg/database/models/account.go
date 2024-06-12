package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type Account struct {
	gorm.Model

	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID string    `gorm:"index:idx_user_id"`

	goserver.AccountNoId
}

func (a *Account) FromDB() goserver.Account {
	return goserver.Account{
		Id:          a.ID.String(),
		Name:        a.Name,
		Type:        a.Type,
		Description: a.Description,
	}
}
