package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	Login          string `gorm:"type:string;primaryKey"`
	StartDate      time.Time
	HashedPassword string
}
