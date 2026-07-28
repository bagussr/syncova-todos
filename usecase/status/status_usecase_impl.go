package usecase

import (
	"context"
	"syncova-todo/delivery/dto"
	"syncova-todo/domain/models"
	repository "syncova-todo/repository/status"
	"time"
)

type statusUsecase struct {
	repository repository.StatusRepository
	timeout    time.Duration
}

func NewStatusUsecase(repository repository.StatusRepository, timeout time.Duration) StatusUsecase {
	return &statusUsecase{
		repository: repository,
		timeout:    timeout,
	}
}

func (s *statusUsecase) GetStatusByProjectID(ctx context.Context, projectId string) ([]models.TodosStatus, error) {
	statuses, err := s.repository.GetStatusByProjectID(ctx, projectId)
	if err != nil {
		return nil, err
	}

	return statuses, nil
}

func (s *statusUsecase) GetStatusByUUID(ctx context.Context, uid string) (models.TodosStatus, error) {
	status, err := s.repository.GetStatusByUUID(ctx, uid)
	if err != nil {
		return models.TodosStatus{}, err
	}

	return status, nil
}

func (s *statusUsecase) CreateStatus(ctx context.Context, status dto.CreateStatusRequest) (models.TodosStatus, error) {
	result, err := s.repository.CreateStatus(ctx, status)
	if err != nil {
		return models.TodosStatus{}, err
	}

	return result, nil
}

func (s *statusUsecase) DeleteStatusByUUID(ctx context.Context, uid string) error {
	err := s.repository.DeleteStatusByUUID(ctx, uid)
	if err != nil {
		return err
	}

	return nil
}
