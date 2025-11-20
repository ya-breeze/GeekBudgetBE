package common

type ContextKey string

const (
	UserIDKey      ContextKey = "userID"
	RequestBodyKey ContextKey = "requestBody"
)
