package repository

import (
	"context"
	"syncova-todo/delivery/dto"
	"syncova-todo/domain/models"
	"syncova-todo/infrastructure/database"
)

type statusRepository struct {
	db *database.PostgresDB
}

var _ StatusRepository = (*statusRepository)(nil)

func NewStatusRepository(db *database.PostgresDB) StatusRepository {
	return &statusRepository{db: db}
}

func (r *statusRepository) GetStatusByProjectID(ctx context.Context, projectId string) ([]models.TodosStatus, error) {
	var status []models.TodosStatus
	var project models.Project
	projectResult, err := r.db.DB.WithContext(ctx).Where("uuid = ?", projectId).First(&project).Rows()
	if err != nil {
		return nil, err
	}
	defer projectResult.Close()

	result, err := r.db.DB.WithContext(ctx).Where("project_id = ?", project.Id).Find(&status).Rows()
	if err != nil {
		return nil, err
	}
	defer result.Close()

	return status, nil
}

func (r *statusRepository) GetStatusByUUID(ctx context.Context, uid string) (models.TodosStatus, error) {
	var status models.TodosStatus

	result, err := r.db.DB.WithContext(ctx).Where("uuid = ?", uid).First(&status).Rows()
	if err != nil {
		return models.TodosStatus{}, err
	}
	defer result.Close()

	return status, nil
}

func (r *statusRepository) CreateStatus(ctx context.Context, status dto.CreateStatusRequest) (models.TodosStatus, error) {
	var project models.Project
	projectResult, err := r.db.DB.WithContext(ctx).Where("uuid = ?", status.ProjectUUID).First(&project).Rows()

	if err != nil {
		return models.TodosStatus{}, err
	}
	defer projectResult.Close()

	todosStatus := models.TodosStatus{
		Status:    status.Status,
		ProjectId: project.Id,
	}

	result := r.db.DB.WithContext(ctx).Create(&todosStatus)
	if result.Error != nil {
		return models.TodosStatus{}, result.Error
	}

	return todosStatus, nil
}

func (r *statusRepository) DeleteStatusByUUID(ctx context.Context, uid string) error {
	result := r.db.DB.WithContext(ctx).Where("uuid = ?", uid).Delete(&models.TodosStatus{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
