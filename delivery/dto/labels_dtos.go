package dto

import (
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/models"
)

type LabelsDto struct {
	Uuid  string `json:"uuid"`
	Label string `json:"label"`
}

type LabelsResponse struct {
	domain.BaseResponse `json:",inline"`
	Data                LabelsDto `json:"data"`
}

type LabelsListResponse struct {
	domain.BaseResponse `json:",inline"`
	Data                []LabelsDto `json:"data"`
}

type CreateLabelRequest struct {
	Label       string `json:"label" binding:"required"`
	ProjectUUID string `json:"project_uuid" binding:"required"`
}

func toLablelsDto(data LabelsDto) LabelsDto {
	return LabelsDto{
		Uuid:  data.Uuid,
		Label: data.Label,
	}
}

func NewLabelSuccessResponse(data LabelsDto, message string) LabelsResponse {
	return LabelsResponse{
		BaseResponse: domain.BaseResponse{
			Success:    true,
			StatusCode: 200,
			Message:    message,
		},
		Data: data,
	}
}

func NewLabelCreateSuccessResponse(data models.LabelsStatuses, message string) LabelsResponse {
	return LabelsResponse{
		BaseResponse: domain.BaseResponse{
			Success:    true,
			StatusCode: 201,
			Message:    message,
		},
		Data: toLablelsDto(LabelsDto{
			Uuid:  data.Uuid,
			Label: data.Label,
		}),
	}
}

func NewLabelsListSuccessResponse(data []models.LabelsStatuses, message string) LabelsListResponse {
	labelsDto := make([]LabelsDto, 0)
	for _, label := range data {
		labelsDto = append(labelsDto, toLablelsDto(LabelsDto{
			Uuid:  label.Uuid,
			Label: label.Label,
		}))
	}
	return LabelsListResponse{
		BaseResponse: domain.BaseResponse{
			Success:    true,
			StatusCode: 200,
			Message:    message,
		},
		Data: labelsDto,
	}
}
