package constants

import "github.com/shopspring/decimal"

type ContextKey string

const (
	UserIDKey       ContextKey = "userID"
	RequestBodyKey  ContextKey = "requestBody"
	ChangeSourceKey ContextKey = "changeSource"
	ForcedImportKey ContextKey = "forced_import_channel"
)

type ChangeSource string

const (
	ChangeSourceUser   ChangeSource = "user"
	ChangeSourceSystem ChangeSource = "system"
)

// ReconciliationTolerance is the maximum allowed difference between app balance and bank balance
// for the account to be considered reconciled.
var ReconciliationTolerance = decimal.NewFromFloat(0.01)
