package usecase

import (
	"context"
	"syncova-todo/delivery/dto"
	"syncova-todo/domain/models"
)

type StatusUsecase interface {
	GetStatusByProjectID(ctx context.Context, projectId string) ([]models.TodosStatus, error)
	GetStatusByUUID(ctx context.Context, uid string) (models.TodosStatus, error)
	CreateStatus(ctx context.Context, status dto.CreateStatusRequest) (models.TodosStatus, error)
	DeleteStatusByUUID(ctx context.Context, uid string) error
}
