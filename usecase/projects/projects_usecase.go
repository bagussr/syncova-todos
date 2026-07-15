package usecase

import (
	"context"
	"syncova-todo/delivery/dto"
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/models"
)

type ProjectsUsecase interface {
	GetProjects(ctx context.Context, request *domain.BasePaginationRequest) (*domain.BaseListResponse, error)
	GetProjectByUUID(ctx context.Context, uid string) (models.Project, error)
	CreateProject(ctx context.Context, request dto.CreateProjectRequest) (models.Project, error)
	UpdateProject(ctx context.Context, uid string, request dto.UpdateProjectRequest) (models.Project, error)
	DeleteProject(ctx context.Context, uid string) error
}
