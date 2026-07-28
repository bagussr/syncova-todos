package dto

import (
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/models"
)

type StatusDto struct {
	UUID   string `json:"uuid"`
	Status string `json:"status"`
}

type StatusResponse struct {
	domain.BaseResponse `json:",inline"`
	Data                StatusDto `json:"data"`
}

type StatusListResponse struct {
	domain.BaseResponse `json:",inline"`
	Data                []StatusDto `json:"data"`
}

type CreateStatusRequest struct {
	Status      string `json:"status" binding:"required"`
	ProjectUUID string `json:"project_uuid" binding:"required"`
}

func NewStatusSuccessResponse(data models.TodosStatus, message string) StatusResponse {
	return StatusResponse{
		BaseResponse: domain.BaseResponse{
			Success:    true,
			StatusCode: 200,
			Message:    message,
		},
		Data: StatusDto{
			UUID:   data.Uuid,
			Status: data.Status,
		},
	}
}

func NewStatusListSuccessResponse(data []models.TodosStatus, message string) StatusListResponse {
	statusDtos := make([]StatusDto, len(data))
	for i, status := range data {
		statusDtos[i] = StatusDto{
			UUID:   status.Uuid,
			Status: status.Status,
		}
	}

	return StatusListResponse{
		BaseResponse: domain.BaseResponse{
			Success:    true,
			StatusCode: 200,
			Message:    message,
		},
		Data: statusDtos,
	}
}

func NewStatusErrorResponse(message string, statusCode int) StatusResponse {
	return StatusResponse{
		BaseResponse: domain.BaseResponse{
			Success:    false,
			StatusCode: statusCode,
			Message:    message,
		},
	}
}
