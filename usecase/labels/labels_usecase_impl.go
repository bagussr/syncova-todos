package usecase

import (
	"syncova-todo/delivery/dto"
	"syncova-todo/domain/models"
	repository "syncova-todo/repository/labels"
	"time"
)

type labelsUsecase struct {
	repository repository.LabelsRepository
	timeout    time.Duration
}

func NewLabelsUsecase(repository repository.LabelsRepository, timeout time.Duration) LabelsUseCase {
	return &labelsUsecase{
		repository: repository,
		timeout:    timeout,
	}
}

func (l *labelsUsecase) GetLabelsByProjectID(projectID string) ([]models.LabelsStatuses, error) {
	labels, err := l.repository.GetLabelsByProjectID(projectID)
	if err != nil {
		return nil, err
	}

	return labels, nil
}

func (l *labelsUsecase) CreateLabel(label dto.CreateLabelRequest) (models.LabelsStatuses, error) {
	newLabel, err := l.repository.CreateLabel(label)
	if err != nil {
		return models.LabelsStatuses{}, err
	}

	return newLabel, nil
}

func (l *labelsUsecase) DeleteLabelByUuid(labelID string) error {
	err := l.repository.DeleteLabelByUuid(labelID)
	if err != nil {
		return err
	}

	return nil
}
