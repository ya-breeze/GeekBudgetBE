package models

import (
	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type Matcher struct {
	gorm.Model

	Name                 string
	OutputDescription    string
	Amount               float64
	CurrencyRegExp       string
	PartnerNameRegExp    string
	PartnerAccountNumber string
	DescriptionRegExp    string
	ExtraRegExp          string
	OutputMovements      []goserver.Movement `gorm:"serializer:json"`

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (t *Matcher) FromDB() goserver.Matcher {
	return goserver.Matcher{
		Id:                   t.ID.String(),
		Name:                 t.Name,
		OutputDescription:    t.OutputDescription,
		Amount:               t.Amount,
		CurrencyRegExp:       t.CurrencyRegExp,
		PartnerNameRegExp:    t.PartnerNameRegExp,
		PartnerAccountNumber: t.PartnerAccountNumber,
		DescriptionRegExp:    t.DescriptionRegExp,
		ExtraRegExp:          t.ExtraRegExp,
		OutputMovements:      t.OutputMovements,
	}
}

func MatcherToDB(m *goserver.MatcherNoId, userID string) *Matcher {
	return &Matcher{
		UserID:               userID,
		Name:                 m.Name,
		OutputDescription:    m.OutputDescription,
		Amount:               m.Amount,
		CurrencyRegExp:       m.CurrencyRegExp,
		PartnerNameRegExp:    m.PartnerNameRegExp,
		PartnerAccountNumber: m.PartnerAccountNumber,
		DescriptionRegExp:    m.DescriptionRegExp,
		ExtraRegExp:          m.ExtraRegExp,
		OutputMovements:      m.OutputMovements,
	}
}
