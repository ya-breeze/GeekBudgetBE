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
	PlaceRegExp                string
	ConfirmationHistory        []bool `gorm:"serializer:json"`
	Image                      string
	Simplified                 bool
	Keywords                   []string `gorm:"serializer:json"`

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (m *Matcher) FromDB() goserver.Matcher {
	count := int32(0)
	for _, v := range m.ConfirmationHistory {
		if v {
			count++
		}
	}
	total := int32(len(m.ConfirmationHistory))

	return goserver.Matcher{
		Id:                         m.ID.String(),
		OutputDescription:          m.OutputDescription,
		OutputAccountId:            m.OutputAccountID,
		OutputTags:                 m.OutputTags,
		CurrencyRegExp:             m.CurrencyRegExp,
		PartnerNameRegExp:          m.PartnerNameRegExp,
		PartnerAccountNumberRegExp: m.PartnerAccountNumberRegExp,
		DescriptionRegExp:          m.DescriptionRegExp,
		ExtraRegExp:                m.ExtraRegExp,
		PlaceRegExp:                m.PlaceRegExp,
		ConfirmationHistory:        m.ConfirmationHistory,
		ConfirmationsCount:         count,
		ConfirmationsTotal:         total,
		Image:                      m.Image,
		Simplified:                 m.Simplified,
		Keywords:                   m.Keywords,
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
		OutputDescription:          m.GetOutputDescription(),
		OutputAccountID:            m.GetOutputAccountId(),
		OutputTags:                 m.GetOutputTags(),
		CurrencyRegExp:             m.GetCurrencyRegExp(),
		PartnerNameRegExp:          m.GetPartnerNameRegExp(),
		PartnerAccountNumberRegExp: m.GetPartnerAccountNumberRegExp(),
		DescriptionRegExp:          m.GetDescriptionRegExp(),
		ExtraRegExp:                m.GetExtraRegExp(),
		PlaceRegExp:                m.GetPlaceRegExp(),
		ConfirmationHistory:        history,
		Image:                      m.GetImage(),
		Simplified:                 m.GetSimplified(),
		Keywords:                   m.GetKeywords(),
	}
}

func MatcherWithoutID(matcher *goserver.Matcher) *goserver.MatcherNoId {
	return &goserver.MatcherNoId{
		OutputDescription:          matcher.OutputDescription,
		OutputAccountId:            matcher.OutputAccountId,
		OutputTags:                 matcher.OutputTags,
		CurrencyRegExp:             matcher.CurrencyRegExp,
		PartnerNameRegExp:          matcher.PartnerNameRegExp,
		PartnerAccountNumberRegExp: matcher.PartnerAccountNumberRegExp,
		DescriptionRegExp:          matcher.DescriptionRegExp,
		ExtraRegExp:                matcher.ExtraRegExp,
		PlaceRegExp:                matcher.PlaceRegExp,
		ConfirmationHistory:        matcher.ConfirmationHistory,
		Image:                      matcher.Image,
		Simplified:                 matcher.Simplified,
		Keywords:                   matcher.Keywords,
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

// IsPerfectMatch returns true if the matcher qualifies as a "perfect match" for auto-conversion.
// A perfect match requires: (1) at least 10 confirmations in history, and
// (2) 100% success rate (all confirmations are true).
func (m *Matcher) IsPerfectMatch() bool {
	if len(m.ConfirmationHistory) < 10 {
		return false
	}

	for _, confirmed := range m.ConfirmationHistory {
		if !confirmed {
			return false
		}
	}

	return true
}
