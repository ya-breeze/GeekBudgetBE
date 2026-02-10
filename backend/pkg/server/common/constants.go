package common

type ContextKey string

const (
	UserIDKey      ContextKey = "userID"
	RequestBodyKey ContextKey = "requestBody"

	// ReconciliationTolerance is the maximum allowed difference between app balance and bank balance
	// for the account to be considered reconciled.
	ReconciliationTolerance = 0.01
)
