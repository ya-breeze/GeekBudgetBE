package models

import (
	"github.com/google/uuid"
	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
	"gorm.io/gorm"
)

type Matcher struct {
	gorm.Model

	Name                       string
	OutputDescription          string
	OutputAccountID            string
	OutputTags                 []string `gorm:"serializer:json"`
	CurrencyRegExp             string
	PartnerNameRegExp          string
	PartnerAccountNumberRegExp string
	DescriptionRegExp          string
	ExtraRegExp                string
	ConfirmationHistory        []bool `gorm:"serializer:json"`

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (t *Matcher) FromDB() goserver.Matcher {
	return goserver.Matcher{
		Id:                         t.ID.String(),
		Name:                       t.Name,
		OutputDescription:          t.OutputDescription,
		OutputAccountId:            t.OutputAccountID,
		OutputTags:                 t.OutputTags,
		CurrencyRegExp:             t.CurrencyRegExp,
		PartnerNameRegExp:          t.PartnerNameRegExp,
		PartnerAccountNumberRegExp: t.PartnerAccountNumberRegExp,
		DescriptionRegExp:          t.DescriptionRegExp,
		ExtraRegExp:                t.ExtraRegExp,
		ConfirmationHistory:        t.ConfirmationHistory,
	}
}

func MatcherToDB(m goserver.MatcherNoIdInterface, userID string) *Matcher {
	// Preserve confirmation history from the incoming model. Ensure non-nil slice.
	history := m.GetConfirmationHistory()
	if history == nil {
		history = make([]bool, 0)
	}

	return &Matcher{
		UserID:                     userID,
		Name:                       m.GetName(),
		OutputDescription:          m.GetOutputDescription(),
		OutputAccountID:            m.GetOutputAccountId(),
		OutputTags:                 m.GetOutputTags(),
		CurrencyRegExp:             m.GetCurrencyRegExp(),
		PartnerNameRegExp:          m.GetPartnerNameRegExp(),
		PartnerAccountNumberRegExp: m.GetPartnerAccountNumberRegExp(),
		DescriptionRegExp:          m.GetDescriptionRegExp(),
		ExtraRegExp:                m.GetExtraRegExp(),
		ConfirmationHistory:        history,
	}
}

func MatcherWithoutID(matcher *goserver.Matcher) *goserver.MatcherNoId {
	return &goserver.MatcherNoId{
		Name:                       matcher.Name,
		OutputDescription:          matcher.OutputDescription,
		OutputAccountId:            matcher.OutputAccountId,
		OutputTags:                 matcher.OutputTags,
		CurrencyRegExp:             matcher.CurrencyRegExp,
		PartnerNameRegExp:          matcher.PartnerNameRegExp,
		PartnerAccountNumberRegExp: matcher.PartnerAccountNumberRegExp,
		DescriptionRegExp:          matcher.DescriptionRegExp,
		ExtraRegExp:                matcher.ExtraRegExp,
		ConfirmationHistory:        matcher.ConfirmationHistory,
	}
}

// GetConfirmationPercentage calculates the percentage of confirmed matches
func (m *Matcher) GetConfirmationPercentage() float64 {
	if len(m.ConfirmationHistory) == 0 {
		return 0.0
	}

	confirmed := 0
	for _, isConfirmed := range m.ConfirmationHistory {
		if isConfirmed {
			confirmed++
		}
	}

	return float64(confirmed) / float64(len(m.ConfirmationHistory)) * 100.0
}

// AddConfirmation adds a new confirmation to the history, maintaining the maximum length
func (m *Matcher) AddConfirmation(confirmed bool, maxLength int) {
	m.ConfirmationHistory = append(m.ConfirmationHistory, confirmed)

	// Maintain maximum length by removing oldest entries
	if len(m.ConfirmationHistory) > maxLength {
		m.ConfirmationHistory = m.ConfirmationHistory[len(m.ConfirmationHistory)-maxLength:]
	}
}

// GetConfirmationHistoryLength returns the current length of confirmation history
func (m *Matcher) GetConfirmationHistoryLength() int {
	return len(m.ConfirmationHistory)
}
