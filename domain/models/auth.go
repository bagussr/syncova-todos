package models

import domain "syncova-todo/domain/base"

type AuthResponse struct {
	domain.BaseResponse
	Data bool `json:"data"`
}

type AuthTokenClaims struct {
	Uuid     string `json:"uuid"`
	Email    string `json:"email"`
	IsActive bool   `json:"is_active"`
}
