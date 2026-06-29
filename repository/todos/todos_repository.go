package repository

import "context"

type TodosRepository interface {
	GetTodos(ctx context.Context) (string, error)
}
