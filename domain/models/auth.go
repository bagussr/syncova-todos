package models

import domain "syncova-todo/domain/base"

type AuthResponse struct {
	domain.BaseResponse
	Data bool `json:"data"`
}
