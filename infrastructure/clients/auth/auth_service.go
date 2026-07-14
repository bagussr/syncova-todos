package clients

import (
	"context"
)

type AuthServiceClient interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
}

type ContextKey string

const (
	ContextKeyUserId ContextKey = "user_id"
	ContextKeyEmail  ContextKey = "email"
	ContextKeyClaims ContextKey = "claims"
)
