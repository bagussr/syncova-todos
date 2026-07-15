package clients

import (
	"context"
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/models"

	"github.com/golang-jwt/jwt/v5"
)

type AuthLocal struct {
	apiKey    string
	secretKey string
}

func NewAuthLocal(apiKey string, secretKey string) *AuthLocal {
	return &AuthLocal{
		apiKey:    apiKey,
		secretKey: secretKey,
	}
}

func (a *AuthLocal) ValidateToken(ctx context.Context, token string) (bool, error) {
	if a.apiKey == "" {
		return false, domain.ErrInternal
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return []byte(a.secretKey), nil
	})

	if err != nil || !parsedToken.Valid {
		return false, domain.ErrUnauthorized
	}

	return true, nil
}

func (a *AuthLocal) ValidateTokenWithClaims(ctx context.Context, token string) (*models.AuthTokenClaims, error) {
	if a.secretKey == "" {
		return nil, domain.ErrInternal
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return []byte(a.secretKey), nil
	})

	if err != nil || !parsedToken.Valid {
		return nil, domain.ErrUnauthorized
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrUnauthorized
	}

	authClaims := &models.AuthTokenClaims{
		Uuid:     claims["uuid"].(string),
		Email:    claims["email"].(string),
		IsActive: claims["is_active"].(bool),
	}

	return authClaims, nil
}
