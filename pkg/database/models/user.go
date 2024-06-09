package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model

	ID             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Email          string
	StartDate      time.Time
	HashedPassword string
}
