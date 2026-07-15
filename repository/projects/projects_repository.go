package repository

import (
	"context"
	"syncova-todo/delivery/dto"
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/models"
)

type ProjectsRepository interface {
	GetProjects(ctx context.Context, request *domain.BasePaginationRequest, userID string) (*domain.BaseListResponse, error)
	GetProjectByUUID(ctx context.Context, uid string) (models.Project, error)
	CreateProject(ctx context.Context, request dto.CreateProjectRequest, userID string) (models.Project, error)
	UpdateProject(ctx context.Context, uid string, request dto.UpdateProjectRequest) (models.Project, error)
	DeleteProject(ctx context.Context, uid string) error
}
