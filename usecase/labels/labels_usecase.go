package usecase

import (
	"syncova-todo/delivery/dto"
	"syncova-todo/domain/models"
)

type LabelsUseCase interface {
	GetLabelsByProjectID(projectID string) ([]models.LabelsStatuses, error)
	CreateLabel(label dto.CreateLabelRequest) (models.LabelsStatuses, error)
	DeleteLabelByUuid(labelID string) error
}
