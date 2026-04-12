package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TransactionDuplicate tracks duplicate relationships between transactions.
// It is stored bidirectionally: if T1 is a duplicate of T2, we store (T1, T2) and (T2, T1).
type TransactionDuplicate struct {
	gorm.Model
	FamilyID       uuid.UUID `gorm:"type:uuid;not null;index:idx_family_trans1_trans2,priority:1"`
	TransactionID1 uuid.UUID `gorm:"index:idx_family_trans1_trans2,priority:2;type:uuid"`
	TransactionID2 uuid.UUID `gorm:"index:idx_family_trans1_trans2,priority:3;type:uuid"`
}
