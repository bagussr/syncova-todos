package repository

import (
	"syncova-todo/delivery/dto"
	"syncova-todo/domain/models"
	"syncova-todo/infrastructure/database"
)

type labelsRepository struct {
	db *database.PostgresDB
}

var _ LabelsRepository = (*labelsRepository)(nil)

var _ LabelsRepository = (*labelsRepository)(nil)

func NewLabelsRepository(db *database.PostgresDB) LabelsRepository {
	return &labelsRepository{db: db}
}

func (r *labelsRepository) GetLabelsByProjectID(projectID string) ([]models.LabelsStatuses, error) {
	var project models.Project
	resultProject, err := r.db.WithContext(r.db.DB.Statement.Context).Where("uuid = ?", projectID).First(&project).Rows()
	if err != nil {
		return nil, err
	}
	defer resultProject.Close()

	var labels []models.LabelsStatuses
	resultLabels, err := r.db.WithContext(r.db.DB.Statement.Context).Where("project_id = ?", project.Id).Find(&labels).Rows()
	if err != nil {
		return nil, err
	}
	defer resultLabels.Close()

	return labels, nil
}

func (r *labelsRepository) CreateLabel(label dto.CreateLabelRequest) (models.LabelsStatuses, error) {
	var project models.Project
	resultProject, err := r.db.WithContext(r.db.DB.Statement.Context).Where("uuid = ?", label.ProjectUUID).First(&project).Rows()
	if err != nil {
		return models.LabelsStatuses{}, err
	}
	defer resultProject.Close()

	newLabel := models.LabelsStatuses{
		Label:     label.Label,
		ProjectId: project.Id,
	}
	result, err := r.db.WithContext(r.db.DB.Statement.Context).Create(&newLabel).Rows()
	if err != nil {
		return models.LabelsStatuses{}, err
	}
	defer result.Close()

	return newLabel, nil
}

func (r *labelsRepository) DeleteLabelByUuid(labelID string) error {
	var label models.LabelsStatuses
	resultLabel, err := r.db.WithContext(r.db.DB.Statement.Context).Where("uuid = ?", labelID).First(&label).Rows()
	if err != nil {
		return err
	}
	defer resultLabel.Close()

	result, err := r.db.WithContext(r.db.DB.Statement.Context).Delete(&label).Rows()
	if err != nil {
		return err
	}
	defer result.Close()

	return nil
}
