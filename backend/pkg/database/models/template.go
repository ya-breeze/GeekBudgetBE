package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type TransactionTemplate struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string
	Description string
	Place       string
	Tags        []string `gorm:"serializer:json"`
	PartnerName string
	Extra       string
	Movements   []goserver.Movement `gorm:"serializer:json"`
	UserID      string              `gorm:"index"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (t *TransactionTemplate) FromDB() goserver.TransactionTemplate {
	return goserver.TransactionTemplate{
		Id:          t.ID.String(),
		Name:        t.Name,
		Description: t.Description,
		Place:       t.Place,
		Tags:        t.Tags,
		PartnerName: t.PartnerName,
		Extra:       t.Extra,
		Movements:   t.Movements,
	}
}

func TemplateToDB(t *goserver.TransactionTemplateNoId, userID string) *TransactionTemplate {
	tags := t.Tags
	if tags == nil {
		tags = make([]string, 0)
	}

	movements := t.Movements
	if movements == nil {
		movements = make([]goserver.Movement, 0)
	}

	return &TransactionTemplate{
		UserID:      userID,
		Name:        t.Name,
		Description: t.Description,
		Place:       t.Place,
		Tags:        tags,
		PartnerName: t.PartnerName,
		Extra:       t.Extra,
		Movements:   movements,
	}
}
