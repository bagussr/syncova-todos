package clients

import (
	"context"
	"syncova-todo/domain/models"
)

type AuthServiceClient interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
}

type AuthServiceLocal interface {
	ValidateToken(ctx context.Context, token string) (bool, error)
	ValidateTokenWithClaims(ctx context.Context, token string) (*models.AuthTokenClaims, error)
}

type ContextKey string

const (
	ContextKeyUserId ContextKey = "user_id"
	ContextKeyEmail  ContextKey = "email"
	ContextKeyClaims ContextKey = "claims"
)
