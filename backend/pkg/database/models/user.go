package models

import (
	"time"

	coremodels "github.com/ya-breeze/kin-core/models"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type User struct {
	coremodels.User
	StartDate          time.Time
	FavoriteCurrencyID string
}

func (u User) FromDB() goserver.User {
	return goserver.User{
		Email:              u.Username,
		StartDate:          u.StartDate,
		FavoriteCurrencyId: u.FavoriteCurrencyID,
	}
}
