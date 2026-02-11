package common

import "github.com/shopspring/decimal"

type ContextKey string

const (
	UserIDKey      ContextKey = "userID"
	RequestBodyKey ContextKey = "requestBody"
)

// ReconciliationTolerance is the maximum allowed difference between app balance and bank balance
// for the account to be considered reconciled.
var ReconciliationTolerance = decimal.NewFromFloat(0.01)
