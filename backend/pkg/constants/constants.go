package constants

import (
	"context"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ContextKey string

const (
	UserIDKey       ContextKey = "userID"
	FamilyIDKey     ContextKey = "familyID"
	RequestBodyKey  ContextKey = "requestBody"
	ChangeSourceKey ContextKey = "changeSource"
	ForcedImportKey ContextKey = "forced_import_channel"
)

// GetFamilyID extracts the family UUID from a request context.
// Returns the UUID and true if found, or zero UUID and false if not.
func GetFamilyID(ctx context.Context) (uuid.UUID, bool) {
	v := ctx.Value(FamilyIDKey)
	id, ok := v.(uuid.UUID)
	return id, ok
}

type ChangeSource string

const (
	ChangeSourceUser   ChangeSource = "user"
	ChangeSourceSystem ChangeSource = "system"
)

// ReconciliationTolerance is the maximum allowed difference between app balance and bank balance
// for the account to be considered reconciled.
var ReconciliationTolerance = decimal.NewFromFloat(0.01)
