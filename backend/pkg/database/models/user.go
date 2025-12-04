package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	StartDate          time.Time
	Login              string `gorm:"unique"`
	HashedPassword     string
	FavoriteCurrencyID string
}

func (u User) FromDB() goserver.User {
	return goserver.User{
		Email:              u.Login,
		StartDate:          u.StartDate,
		FavoriteCurrencyId: u.FavoriteCurrencyID,
	}
}
