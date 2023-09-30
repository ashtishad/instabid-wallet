package domain

import "time"

const (
	UserCredentialEmail    = "email"
	UserCredentialUsername = "username"

	AccessTokenDuration = time.Hour
	TokenTypeAccess     = "access_token"
)

type ContextKey string

const (
	UserCredentialKey ContextKey = "credential"
)
