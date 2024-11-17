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
	OutputTags           []string `gorm:"serializer:json"`
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
		OutputTags:           t.OutputTags,
		CurrencyRegExp:       t.CurrencyRegExp,
		PartnerNameRegExp:    t.PartnerNameRegExp,
		PartnerAccountNumber: t.PartnerAccountNumber,
		DescriptionRegExp:    t.DescriptionRegExp,
		ExtraRegExp:          t.ExtraRegExp,
	}
}

func MatcherToDB(m goserver.MatcherNoIdInterface, userID string) *Matcher {
	return &Matcher{
		UserID:               userID,
		Name:                 m.GetName(),
		OutputDescription:    m.GetOutputDescription(),
		OutputAccountID:      m.GetOutputAccountId(),
		OutputTags:           m.GetOutputTags(),
		CurrencyRegExp:       m.GetCurrencyRegExp(),
		PartnerNameRegExp:    m.GetPartnerNameRegExp(),
		PartnerAccountNumber: m.GetPartnerAccountNumber(),
		DescriptionRegExp:    m.GetDescriptionRegExp(),
		ExtraRegExp:          m.GetExtraRegExp(),
	}
}

func MatcherWithoutID(matcher *goserver.Matcher) *goserver.MatcherNoId {
	return &goserver.MatcherNoId{
		Name:                 matcher.Name,
		OutputDescription:    matcher.OutputDescription,
		OutputAccountId:      matcher.OutputAccountId,
		OutputTags:           matcher.OutputTags,
		CurrencyRegExp:       matcher.CurrencyRegExp,
		PartnerNameRegExp:    matcher.PartnerNameRegExp,
		PartnerAccountNumber: matcher.PartnerAccountNumber,
		DescriptionRegExp:    matcher.DescriptionRegExp,
		ExtraRegExp:          matcher.ExtraRegExp,
	}
}
