package usecase

import (
	"context"
	"syncova-todo/delivery/dto"
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/models"
	repository "syncova-todo/repository/todos"
	"time"
)

type todosUsecase struct {
	repository repository.TodosRepository
	timeout    time.Duration
}

func NewTodosUsecase(repository repository.TodosRepository, timeout time.Duration) TodosUsecase {
	return &todosUsecase{
		repository: repository,
		timeout:    timeout,
	}
}

func (u *todosUsecase) GetTodos(ctx context.Context, request *domain.BasePaginationRequest, userID string) (*domain.BaseListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	todos, err := u.repository.GetTodos(ctx, request, userID)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func (u *todosUsecase) GetTodoByUuid(ctx context.Context, uuid string) (models.Todos, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	todos, err := u.repository.GetTodoByUuid(ctx, uuid)
	if err != nil {
		return models.Todos{}, err
	}

	return todos, nil
}

func (u *todosUsecase) CreateTodo(ctx context.Context, data dto.CreateTodoRequest, userID string) (models.Todos, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repository.CreateTodo(ctx, data, userID)
}

func (u *todosUsecase) UpdateTodo(ctx context.Context, uuid string, data dto.UpdateTodoRequest) (models.Todos, error) {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repository.UpdateTodo(ctx, uuid, data)
}

func (u *todosUsecase) DeleteTodo(ctx context.Context, uuid string) error {
	ctx, cancel := context.WithTimeout(ctx, u.timeout)
	defer cancel()

	return u.repository.DeleteTodo(ctx, uuid)
}
