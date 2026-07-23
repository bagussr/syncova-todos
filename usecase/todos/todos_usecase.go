package usecase

import (
	"context"
	"syncova-todo/delivery/dto"
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/models"
)

type TodosUsecase interface {
	GetTodos(ctx context.Context, request *domain.BasePaginationRequest, userID string) (*domain.BaseListResponse, error)
	GetTodoByUuid(ctx context.Context, uuid string) (models.Todos, error)
	CreateTodo(ctx context.Context, data dto.CreateTodoRequest, userID string) (models.Todos, error)
	UpdateTodo(ctx context.Context, uuid string, data dto.UpdateTodoRequest) (models.Todos, error)
	DeleteTodo(ctx context.Context, uuid string) error
}
