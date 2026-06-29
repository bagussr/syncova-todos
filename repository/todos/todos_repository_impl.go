package repository

import "context"

type todosRepository struct {
}

var _ TodosRepository = (*todosRepository)(nil)

func NewTodosRepository() TodosRepository {
	return &todosRepository{}
}

func (r *todosRepository) GetTodos(ctx context.Context) (string, error) {

	return "Sample Todos", nil
}
