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
	OutputAccountID      string
	CurrencyRegExp       string
	PartnerNameRegExp    string
	PartnerAccountNumber string
	DescriptionRegExp    string
	ExtraRegExp          string

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (t *Matcher) FromDB() goserver.Matcher {
	return goserver.Matcher{
		Id:                   t.ID.String(),
		Name:                 t.Name,
		OutputDescription:    t.OutputDescription,
		OutputAccountId:      t.OutputAccountID,
		CurrencyRegExp:       t.CurrencyRegExp,
		PartnerNameRegExp:    t.PartnerNameRegExp,
		PartnerAccountNumber: t.PartnerAccountNumber,
		DescriptionRegExp:    t.DescriptionRegExp,
		ExtraRegExp:          t.ExtraRegExp,
	}
}

func MatcherToDB(m *goserver.MatcherNoId, userID string) *Matcher {
	return &Matcher{
		UserID:               userID,
		Name:                 m.Name,
		OutputDescription:    m.OutputDescription,
		OutputAccountID:      m.OutputAccountId,
		CurrencyRegExp:       m.CurrencyRegExp,
		PartnerNameRegExp:    m.PartnerNameRegExp,
		PartnerAccountNumber: m.PartnerAccountNumber,
		DescriptionRegExp:    m.DescriptionRegExp,
		ExtraRegExp:          m.ExtraRegExp,
	}
}

func MatcherWithoutID(matcher *goserver.Matcher) *goserver.MatcherNoId {
	return &goserver.MatcherNoId{
		Name:                 matcher.Name,
		OutputDescription:    matcher.OutputDescription,
		OutputAccountId:      matcher.OutputAccountId,
		CurrencyRegExp:       matcher.CurrencyRegExp,
		PartnerNameRegExp:    matcher.PartnerNameRegExp,
		PartnerAccountNumber: matcher.PartnerAccountNumber,
		DescriptionRegExp:    matcher.DescriptionRegExp,
		ExtraRegExp:          matcher.ExtraRegExp,
	}
}
