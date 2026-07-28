package dto

import (
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/models"
	"time"
)

type ProjectDto struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date"`
}

type ProjectResponse struct {
	domain.BaseResponse `json:",inline"`
	Data                ProjectDto `json:"data"`
}

type ProjectsListResponse struct {
	domain.BaseListResponse
	Data interface{} `json:"data"`
}

type CreateProjectRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date" binding:"required"`
}

type UpdateProjectRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date"`
}

func toProjectDto(data models.Project) ProjectDto {
	return ProjectDto{
		UUID:        data.Uuid,
		Name:        data.Name,
		Description: data.Description,
		Status:      string(data.Status),
		DueDate:     data.DueDate,
	}
}

func NewProjectSuccessResponse(data models.Project, message string) ProjectResponse {
	return ProjectResponse{
		BaseResponse: domain.BaseResponse{
			Success:    true,
			StatusCode: 200,
			Message:    message,
		},
		Data: ProjectDto{
			UUID:        data.Uuid,
			Name:        data.Name,
			Description: data.Description,
			Status:      string(data.Status),
			DueDate:     data.DueDate,
		},
	}
}

func NewProjectsListSuccessResponse(data *domain.BaseListResponse, message string) ProjectsListResponse {
	projects := []ProjectDto{}

	if items, ok := data.Data.([]models.Project); ok {
		for _, item := range items {
			projects = append(projects, toProjectDto(item))
		}
	} else if items, ok := data.Data.(*[]models.Project); ok && items != nil {
		for _, item := range *items {
			projects = append(projects, toProjectDto(item))
		}
	}

	return ProjectsListResponse{
		BaseListResponse: domain.BaseListResponse{
			Success:    true,
			StatusCode: 200,
			Message:    message,
			Page:       data.Page,
			PerPage:    data.PerPage,
			Total:      data.Total,
		},
		Data: projects,
	}
}

func NewProjectCreateSuccessResponse(data models.Project, message string) ProjectResponse {
	return ProjectResponse{
		BaseResponse: domain.BaseResponse{
			Success:    true,
			StatusCode: 201,
			Message:    message,
		},
		Data: ProjectDto{
			UUID:        data.Uuid,
			Name:        data.Name,
			Description: data.Description,
			Status:      string(data.Status),
			DueDate:     data.DueDate,
		},
	}
}

func NewProjectErrorResponse(data models.Project, message string) ProjectResponse {
	return ProjectResponse{
		BaseResponse: domain.BaseResponse{
			Success:    false,
			StatusCode: 500,
			Message:    message,
		},
		Data: ProjectDto{
			UUID:        data.Uuid,
			Name:        data.Name,
			Description: data.Description,
			Status:      string(data.Status),
			DueDate:     data.DueDate,
		},
	}
}
