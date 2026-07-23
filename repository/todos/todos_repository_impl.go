package repository

import (
	"context"

	"syncova-todo/delivery/dto"
	domain "syncova-todo/domain/base"
	"syncova-todo/domain/enums"
	"syncova-todo/domain/models"
	"syncova-todo/infrastructure/database"
)

type todosRepository struct {
	db *database.PostgresDB
}

var _ TodosRepository = (*todosRepository)(nil)

func NewTodosRepository(db *database.PostgresDB) TodosRepository {
	return &todosRepository{db: db}
}

func (r *todosRepository) GetTodos(ctx context.Context, request *domain.BasePaginationRequest, userID string) (*domain.BaseListResponse, error) {
	var todos []models.Todos

	result, err := r.db.Paginated(ctx, request, &todos, []string{"title"}, userID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *todosRepository) GetTodoByUuid(ctx context.Context, uuid string) (models.Todos, error) {
	var todo models.Todos

	result := r.db.DB.WithContext(ctx).Where("uuid = ?", uuid).First(&todo)
	if result.Error != nil {
		return models.Todos{}, result.Error
	}

	return todo, nil
}

func (r *todosRepository) CreateTodo(ctx context.Context, data dto.CreateTodoRequest, userID string) (models.Todos, error) {
	todo := models.Todos{
		Title:       data.Title,
		Description: data.Description,
		DueDate:     data.DueDate,
		StartDate:   data.StartDate,
		Priority:    enums.Priority(data.Priority),
		UserId:      userID,
	}

	result := r.db.DB.WithContext(ctx).Create(&todo)
	if result.Error != nil {
		return models.Todos{}, result.Error
	}

	return todo, nil
}

func (r *todosRepository) UpdateTodo(ctx context.Context, uuid string, data dto.UpdateTodoRequest) (models.Todos, error) {
	var todo models.Todos

	result := r.db.DB.WithContext(ctx).Where("uuid = ?", uuid).First(&todo)
	if result.Error != nil {
		return models.Todos{}, result.Error
	}

	if data.Title != "" {
		todo.Title = data.Title
	}
	if data.Description != "" {
		todo.Description = data.Description
	}
	if data.DueDate != "" {
		todo.DueDate = data.DueDate
	}
	if data.StartDate != "" {
		todo.StartDate = data.StartDate
	}
	if data.Priority != "" {
		todo.Priority = enums.Priority(data.Priority)
	}

	result = r.db.DB.WithContext(ctx).Save(&todo)
	if result.Error != nil {
		return models.Todos{}, result.Error
	}

	return todo, nil
}

func (r *todosRepository) DeleteTodo(ctx context.Context, uuid string) error {
	result := r.db.DB.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Todos{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
