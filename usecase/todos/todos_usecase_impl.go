package usecase

import (
	"context"
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

func (t *todosUsecase) GetTodos(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, t.timeout)
	defer cancel()

	todos, err := t.repository.GetTodos(ctx)
	if err != nil {
		return "", err
	}

	return todos, nil
}
