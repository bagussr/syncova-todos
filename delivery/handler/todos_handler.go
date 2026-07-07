package handler

import (
	"syncova-todo/delivery/dto"
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

// GetTodos godoc
// @Summary Get all todos
// @Description Get all todos
// @Tags todos
// @Accept json
// @Produce json
// @Success 200 {object} dto.TodoResponse
// @Router /todos [get]
func (h *TodosHandler) GetTodos(ctx *gin.Context) {
	todos, err := h.Usecase.GetTodos(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.NewTodoSuccessResponse(todos, "Get todos successfully")

	ctx.JSON(200, response)
}
