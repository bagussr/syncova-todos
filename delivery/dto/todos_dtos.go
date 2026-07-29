package dto

import (
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/models"
)

type TodosDto struct {
	UUID        string `json:"uuid"`
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	StartDate   string `json:"start_date"`
	Priority    string `json:"priority"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type TodoResponse struct {
	domain.BaseResponse `json:",inline"`
	Data                TodosDto `json:"data"`
}

type TodosListResponse struct {
	domain.BaseListResponse
	Data interface{} `json:"data"`
}

type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
	DueDate     string `json:"due_date" validate:"required"`
	StartDate   string `json:"start_date" validate:"required"`
	Priority    string `json:"priority" validate:"oneof=low medium high" default:"low"`
}

type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	StartDate   string `json:"start_date"`
	Priority    string `json:"priority"`
}

func toTodosDto(data models.Todos) TodosDto {
	return TodosDto{
		UUID:        data.Uuid,
		Title:       data.Title,
		Description: data.Description,
		DueDate:     data.DueDate,
		StartDate:   data.StartDate,
		Priority:    string(data.Priority),
		CreatedAt:   data.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   data.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func NewTodoSuccessResponse(data models.Todos, message string) TodoResponse {
	return TodoResponse{
		BaseResponse: domain.BaseResponse{
			Success:    true,
			StatusCode: 200,
			Message:    message,
		},
		Data: toTodosDto(data),
	}
}

func NewTodosListSuccessResponse(data *domain.BaseListResponse, message string) TodosListResponse {
	todos := make([]TodosDto, 0)

	if items, ok := data.Data.(*[]models.Todos); ok && items != nil {
		for _, item := range *items {
			todos = append(todos, toTodosDto(item))
		}
	} else if items, ok := data.Data.([]models.Todos); ok {
		for _, item := range items {
			todos = append(todos, toTodosDto(item))
		}
	}

	return TodosListResponse{
		BaseListResponse: domain.BaseListResponse{
			Success:    true,
			StatusCode: 200,
			Message:    message,
			Page:       data.Page,
			PerPage:    data.PerPage,
			Total:      data.Total,
		},
		Data: todos,
	}
}

func NewTodosCreateSuccessResponse(data models.Todos, message string) TodoResponse {
	return TodoResponse{
		BaseResponse: domain.BaseResponse{
			Success:    true,
			StatusCode: 201,
			Message:    message,
		},
		Data: toTodosDto(data),
	}
}

func NewTodoErrorResponse(message string, statusCode int) TodoResponse {
	return TodoResponse{
		BaseResponse: domain.BaseResponse{
			Success:    false,
			StatusCode: statusCode,
			Message:    message,
		},
	}
}
