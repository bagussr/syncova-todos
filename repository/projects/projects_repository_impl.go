package repository

import (
	"context"
	"syncova-todo/delivery/dto"
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/enums"
	"syncova-todo/domain/models"
	"syncova-todo/infrastructure/database"
)

type projectsRepository struct {
	db *database.PostgresDB
}

var _ ProjectsRepository = (*projectsRepository)(nil)

func NewProjectsRepository(db *database.PostgresDB) ProjectsRepository {
	return &projectsRepository{db: db}
}

func (r *projectsRepository) GetProjects(ctx context.Context, request *domain.BasePaginationRequest, userID string) (*domain.BaseListResponse, error) {
	var projects []models.Project

	result, err := r.db.Paginated(ctx, request, &projects, []string{"name"}, userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *projectsRepository) GetProjectByUUID(ctx context.Context, uid string) (models.Project, error) {
	var project models.Project

	result, err := r.db.DB.WithContext(ctx).Where("uuid = ?", uid).First(&project).Rows()
	if err != nil {
		return models.Project{}, err
	}
	defer result.Close()

	return project, nil
}

func (r *projectsRepository) CreateProject(ctx context.Context, request dto.CreateProjectRequest, userID string) (models.Project, error) {

	project := models.Project{
		Name:        request.Name,
		Description: request.Description,
		Status:      enums.NotStarted,
		DueDate:     request.DueDate,
		UserId:      userID,
	}

	result := r.db.DB.WithContext(ctx).Create(&project)
	if result.Error != nil {
		return models.Project{}, result.Error
	}

	return project, nil
}

func (r *projectsRepository) UpdateProject(ctx context.Context, uid string, request dto.UpdateProjectRequest) (models.Project, error) {
	var project models.Project

	result := r.db.DB.WithContext(ctx).Where("uuid = ?", uid).First(&project)
	if result.Error != nil {
		return models.Project{}, result.Error
	}

	if request.Name != "" {
		project.Name = request.Name
	}
	if request.Description != "" {
		project.Description = request.Description
	}
	if !request.DueDate.IsZero() {
		project.DueDate = request.DueDate
	}

	result = r.db.DB.WithContext(ctx).Save(&project)
	if result.Error != nil {
		return models.Project{}, result.Error
	}

	return project, nil
}

func (r *projectsRepository) DeleteProject(ctx context.Context, uid string) error {
	result := r.db.DB.WithContext(ctx).Where("uuid = ?", uid).Delete(&models.Project{})
	return result.Error
}
