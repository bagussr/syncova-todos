package handler

import (
	"syncova-todo/domain"
	usecase "syncova-todo/usecase/todos"

	"github.com/gin-gonic/gin"
)

type TodosHandler struct {
	Usecase usecase.TodosUsecase
}

func NewTodosHandler(usecase usecase.TodosUsecase) *TodosHandler {
	return &TodosHandler{
		Usecase: usecase,
	}
}

// GET /todos
func (h *TodosHandler) GetTodos(ctx *gin.Context) {
	todos, err := h.Usecase.GetTodos(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := domain.NewSuccessResponse(todos, "Get todos successfully")

	ctx.JSON(200, response)
}
