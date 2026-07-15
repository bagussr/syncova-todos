package usecase

import (
	"context"
	"syncova-todo/delivery/dto"
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/models"
	repository "syncova-todo/repository/projects"
	"time"
)

type projectsUsecase struct {
	repository repository.ProjectsRepository
	timeout    time.Duration
}

func NewProjectsUsecase(repository repository.ProjectsRepository, timeout time.Duration) ProjectsUsecase {
	return &projectsUsecase{
		repository: repository,
		timeout:    timeout,
	}
}

func (p *projectsUsecase) GetProjects(ctx context.Context, request *domain.BasePaginationRequest) (*domain.BaseListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	projects, err := p.repository.GetProjects(ctx, request)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (p *projectsUsecase) GetProjectByUUID(ctx context.Context, uid string) (models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	project, err := p.repository.GetProjectByUUID(ctx, uid)
	if err != nil {
		return models.Project{}, err
	}

	return project, nil
}

func (p *projectsUsecase) CreateProject(ctx context.Context, request dto.CreateProjectRequest) (models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	result, err := p.repository.CreateProject(ctx, request)
	if err != nil {
		return models.Project{}, err
	}

	return result, nil
}

func (p *projectsUsecase) UpdateProject(ctx context.Context, uid string, request dto.UpdateProjectRequest) (models.Project, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	result, err := p.repository.UpdateProject(ctx, uid, request)
	if err != nil {
		return models.Project{}, err
	}

	return result, nil
}

func (p *projectsUsecase) DeleteProject(ctx context.Context, uid string) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	err := p.repository.DeleteProject(ctx, uid)
	if err != nil {
		return err
	}

	return nil
}
