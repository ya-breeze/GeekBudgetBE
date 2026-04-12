package models

import (
	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type TransactionTemplate struct {
	gorm.Model

	Name        string
	Description string
	Place       string
	Tags        []string `gorm:"serializer:json"`
	PartnerName string
	Extra       string
	Movements   []goserver.Movement `gorm:"serializer:json"`

	FamilyID uuid.UUID `gorm:"type:uuid;index;not null"`
	ID       uuid.UUID `gorm:"type:uuid;primaryKey"`
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

func TemplateToDB(t goserver.TransactionTemplateNoIdInterface, familyID uuid.UUID) *TransactionTemplate {
	tags := t.GetTags()
	if tags == nil {
		tags = make([]string, 0)
	}

	movements := t.GetMovements()
	if movements == nil {
		movements = make([]goserver.Movement, 0)
	}

	return &TransactionTemplate{
		FamilyID:    familyID,
		Name:        t.GetName(),
		Description: t.GetDescription(),
		Place:       t.GetPlace(),
		Tags:        tags,
		PartnerName: t.GetPartnerName(),
		Extra:       t.GetExtra(),
		Movements:   movements,
	}
}
