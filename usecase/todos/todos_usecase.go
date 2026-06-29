package usecase

import "context"

type TodosUsecase interface {
	GetTodos(ctx context.Context) (string, error)
}
