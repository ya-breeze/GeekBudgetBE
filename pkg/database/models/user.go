package models

import (
	"time"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Login          string `gorm:"type:string;primaryKey"`
	StartDate      time.Time
	HashedPassword string
}

func (u User) FromDB() goserver.User {
	return goserver.User{
		Email:     u.Login,
		StartDate: u.StartDate,
	}
}
