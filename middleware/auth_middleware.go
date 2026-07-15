package middleware

import (
	"net/http"
	"strings"
	domain "syncova-todo/domain/base"
	clients "syncova-todo/infrastructure/clients/auth"

	"github.com/gin-gonic/gin"
)

type AuthMiddlewareConfig struct {
	AuthClient   clients.AuthServiceClient
	AuthLocal    clients.AuthServiceLocal
	ExcludePaths []string
}

// NewAuthMiddleware creates authentication middleware
func NewAuthMiddleware(config AuthMiddlewareConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip auth for excluded paths
		for _, path := range config.ExcludePaths {
			if c.Request.URL.Path == path {
				c.Next()
				return
			}
		}

		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.NewBaseResponse(
				false, 401, "Authorization header required",
			))
			return
		}

		// Validate Bearer format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.NewBaseResponse(
				false, 401, "Invalid authorization format. Use: Bearer <token>",
			))
			return
		}

		token := parts[1]

		// Call auth microservice to validate token with claims
		claims, err := config.AuthLocal.ValidateTokenWithClaims(c.Request.Context(), token)
		if err != nil {
			// Handle specific errors
			if err == domain.ErrUnauthorized {
				c.AbortWithStatusJSON(http.StatusUnauthorized, domain.NewBaseResponse(
					false, 401, "Invalid or expired token",
				))
				return
			}

			if err == domain.ErrInternal {
				c.AbortWithStatusJSON(http.StatusInternalServerError, domain.NewBaseResponse(
					false, 500, "Internal server error",
				))
				return
			}

			// Auth service unavailable
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, domain.NewBaseResponse(
				false, 503, "Authentication service unavailable",
			))
			return
		}

		// Store claims in Gin context for downstream handlers
		c.Set(string(clients.ContextKeyUserId), claims.Uuid)
		c.Set(string(clients.ContextKeyEmail), claims.Email)
		c.Set(string(clients.ContextKeyClaims), claims)

		c.Next()
	}
}

// Helper functions to extract user data from context

// GetUserID retrieves user ID from Gin context
func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get(string(clients.ContextKeyUserId))
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

// GetUserEmail retrieves user email from Gin context
func GetUserEmail(c *gin.Context) (string, bool) {
	email, exists := c.Get(string(clients.ContextKeyEmail))
	if !exists {
		return "", false
	}
	e, ok := email.(string)
	return e, ok
}

// GetTokenClaims retrieves full claims from Gin context
func GetTokenClaims(c *gin.Context) (bool, bool) {
	claims, exists := c.Get(string(clients.ContextKeyClaims))
	if !exists {
		return false, false
	}
	casted, ok := claims.(bool)
	return casted, ok
}

// RequireAuth is a convenience wrapper that applies auth with default config
func RequireAuth(authClient clients.AuthServiceClient, authLocal clients.AuthServiceLocal) gin.HandlerFunc {
	return NewAuthMiddleware(AuthMiddlewareConfig{
		AuthClient: authClient,
		AuthLocal:  authLocal,
	})
}

// RequireRole enforces specific role
func RequireRole(authClient clients.AuthServiceClient, role string) gin.HandlerFunc {
	return NewAuthMiddleware(AuthMiddlewareConfig{
		AuthClient: authClient,
	})
}

// OptionalAuth allows requests without token but populates context if present
func OptionalAuth(authClient clients.AuthServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.Next()
			return
		}

		claims, err := authClient.ValidateToken(c.Request.Context(), parts[1])
		if err != nil {
			c.Next()
			return
		}

		// c.Set(string(clients.ContextKeyUserId), claims.UserId)
		// c.Set(string(clients.ContextKeyEmail), claims.Email)
		c.Set(string(clients.ContextKeyClaims), claims)

		c.Next()
	}
}
