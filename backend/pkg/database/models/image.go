package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Image struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;"`
	Data        []byte    `gorm:"type:blob"`
	ContentType string
}

func (i *Image) BeforeCreate(tx *gorm.DB) (err error) {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	return
}
