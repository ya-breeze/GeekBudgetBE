package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/ya-breeze/geekbudgetbe/pkg/generated/goserver"
)

type Transaction struct {
	gorm.Model

	Date           time.Time
	Description    string
	Place          string
	Tags           []string `gorm:"serializer:json"`
	PartnerName    string
	PartnerAccount string
	// Internal bank's ID to be able to match later if necessary
	PartnerInternalID string
	// Stores extra data about transaction. For example could hold \"variable symbol\"
	// to distinguish payment for the same account, but with different meaning
	Extra string
	// Stores FULL unprocessed transactions which was source of this transaction.
	// Could be used later for detailed analysis
	UnprocessedSources string
	// IDs of unprocessed transaction - to match later
	ExternalIDs []string            `gorm:"serializer:json"`
	Movements   []goserver.Movement `gorm:"serializer:json"`

	// MatcherID is the ID of the matcher used for this conversion (if any)
	MatcherID *uuid.UUID `gorm:"type:uuid"`
	// IsAuto is true if this transaction was converted automatically by the matcher
	IsAuto bool

	SuspiciousReasons []string `gorm:"serializer:json"`

	// MergedIntoID is set when this transaction was marked as duplicate of another
	MergedIntoID *uuid.UUID `gorm:"type:uuid;index"`
	// MergedAt records when this transaction was merged
	MergedAt *time.Time

	UserID string    `gorm:"index"`
	ID     uuid.UUID `gorm:"type:uuid;primaryKey"`
}

func (t *Transaction) FromDB() goserver.Transaction {
	var matcherID string
	if t.MatcherID != nil {
		matcherID = t.MatcherID.String()
	}

	var mergedIntoID string
	if t.MergedIntoID != nil {
		mergedIntoID = t.MergedIntoID.String()
	}

	var mergedAt time.Time
	if t.MergedAt != nil {
		mergedAt = *t.MergedAt
	}

	return goserver.Transaction{
		Id:                 t.ID.String(),
		Date:               t.Date,
		Description:        t.Description,
		Place:              t.Place,
		Tags:               t.Tags,
		PartnerName:        t.PartnerName,
		PartnerAccount:     t.PartnerAccount,
		PartnerInternalId:  t.PartnerInternalID,
		Extra:              t.Extra,
		UnprocessedSources: t.UnprocessedSources,
		ExternalIds:        t.ExternalIDs,
		Movements:          t.Movements,
		MatcherId:          matcherID,
		IsAuto:             t.IsAuto,
		SuspiciousReasons:  t.SuspiciousReasons,
		MergedIntoId:       mergedIntoID,
		MergedAt:           mergedAt,
	}
}

func (t *Transaction) WithoutID() *goserver.TransactionNoId {
	var matcherID string
	if t.MatcherID != nil {
		matcherID = t.MatcherID.String()
	}

	var mergedIntoID string
	if t.MergedIntoID != nil {
		mergedIntoID = t.MergedIntoID.String()
	}

	var mergedAt time.Time
	if t.MergedAt != nil {
		mergedAt = *t.MergedAt
	}

	return &goserver.TransactionNoId{
		Date:               t.Date,
		Description:        t.Description,
		Place:              t.Place,
		Tags:               t.Tags,
		PartnerName:        t.PartnerName,
		PartnerAccount:     t.PartnerAccount,
		PartnerInternalId:  t.PartnerInternalID,
		Extra:              t.Extra,
		UnprocessedSources: t.UnprocessedSources,
		ExternalIds:        t.ExternalIDs,
		Movements:          t.Movements,
		MatcherId:          matcherID,
		IsAuto:             t.IsAuto,
		SuspiciousReasons:  t.SuspiciousReasons,
		MergedIntoId:       mergedIntoID,
		MergedAt:           mergedAt,
	}
}

func TransactionToDB(transaction goserver.TransactionNoIdInterface, userID string) *Transaction {
	var matcherID *uuid.UUID
	if transaction.GetMatcherId() != "" {
		id, err := uuid.Parse(transaction.GetMatcherId())
		if err == nil {
			matcherID = &id
		}
	}

	var mergedIntoID *uuid.UUID
	if transaction.GetMergedIntoId() != "" {
		id, err := uuid.Parse(transaction.GetMergedIntoId())
		if err == nil {
			mergedIntoID = &id
		}
	}

	var mergedAt *time.Time
	if !transaction.GetMergedAt().IsZero() {
		t := transaction.GetMergedAt()
		mergedAt = &t
	}

	return &Transaction{
		Date:               transaction.GetDate(),
		Description:        transaction.GetDescription(),
		Place:              transaction.GetPlace(),
		Tags:               transaction.GetTags(),
		PartnerName:        transaction.GetPartnerName(),
		PartnerAccount:     transaction.GetPartnerAccount(),
		PartnerInternalID:  transaction.GetPartnerInternalId(),
		Extra:              transaction.GetExtra(),
		UnprocessedSources: transaction.GetUnprocessedSources(),
		ExternalIDs:        transaction.GetExternalIds(),
		Movements:          transaction.GetMovements(),
		MatcherID:          matcherID,
		IsAuto:             transaction.GetIsAuto(),
		SuspiciousReasons:  transaction.GetSuspiciousReasons(),
		MergedIntoID:       mergedIntoID,
		MergedAt:           mergedAt,
		UserID:             userID,
	}
}

func TransactionWithoutID(transaction *goserver.Transaction) *goserver.TransactionNoId {
	return &goserver.TransactionNoId{
		Date:               transaction.Date,
		Description:        transaction.Description,
		Place:              transaction.Place,
		Tags:               transaction.Tags,
		PartnerName:        transaction.PartnerName,
		PartnerAccount:     transaction.PartnerAccount,
		PartnerInternalId:  transaction.PartnerInternalId,
		Extra:              transaction.Extra,
		UnprocessedSources: transaction.UnprocessedSources,
		ExternalIds:        transaction.ExternalIds,
		Movements:          transaction.Movements,
		MatcherId:          transaction.MatcherId,
		IsAuto:             transaction.IsAuto,
		SuspiciousReasons:  transaction.SuspiciousReasons,
		MergedIntoId:       transaction.MergedIntoId,
		MergedAt:           transaction.MergedAt,
	}
}

func MovementsToAPI(movements []goserver.Movement) []goserver.Movement {
	if movements == nil {
		return []goserver.Movement{}
	}
	return movements
}
