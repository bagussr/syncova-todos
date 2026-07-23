package handler

import (
	"strconv"
	"syncova-todo/delivery/dto"
	domain "syncova-todo/domain/base"
	"syncova-todo/middleware"
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
// @security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Number of items per page" default(10)
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort query string false "Sort by field" default(desc)
// @Param search query string false "Search by name or description"
// @Success 200 {object} dto.TodoResponse
// @Router /todos [get]
func (h *TodosHandler) GetTodos(ctx *gin.Context) {
	query := ctx.Request.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	perPage, _ := strconv.Atoi(query.Get("per_page"))
	request := &domain.BasePaginationRequest{
		Page:    page,
		PerPage: perPage,
		SortBy:  query.Get("sort_by"),
		Sort:    query.Get("sort"),
		Search:  query.Get("search"),
	}

	userID, _ := middleware.GetUserID(ctx)

	todos, err := h.Usecase.GetTodos(ctx, request, userID)
	if err != nil {
		response := dto.NewTodoErrorResponse(err.Error(), 500)
		ctx.JSON(500, response)
		return
	}

	response := dto.NewTodosListSuccessResponse(todos, "Get todos successfully")
	ctx.JSON(200, response)
}

// GetTodoByUuid godoc
// @Summary Get todo by uuid
// @Description Get todo by uuid
// @Tags todos
// @Accept json
// @Produce json
// @security BearerAuth
// @Param uuid path string true "Todo UUID"
// @Success 200 {object} dto.TodoResponse
// @Router /todos/{uuid} [get]
func (h *TodosHandler) GetTodoByUuid(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	todo, err := h.Usecase.GetTodoByUuid(ctx, uuid)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.NewTodoSuccessResponse(todo, "Get todo successfully")
	ctx.JSON(200, response)
}

// CreateTodo godoc
// @Summary Create todo
// @Description Create todo
// @Tags todos
// @Accept json
// @Produce json
// @security BearerAuth
// @Param request body dto.CreateTodoRequest true "Create todo request"
// @Success 201 {object} dto.TodoResponse
// @Router /todos [post]
func (h *TodosHandler) CreateTodo(ctx *gin.Context) {
	var request dto.CreateTodoRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	userID, _ := middleware.GetUserID(ctx)

	todo, err := h.Usecase.CreateTodo(ctx, request, userID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.NewTodosCreateSuccessResponse(todo, "Create todo successfully")
	ctx.JSON(201, response)
}

// UpdateTodo godoc
// @Summary Update todo
// @Description Update todo
// @Tags todos
// @Accept json
// @Produce json
// @security BearerAuth
// @Param uuid path string true "Todo UUID"
// @Param request body dto.UpdateTodoRequest true "Update todo request"
// @Success 200 {object} dto.TodoResponse
// @Router /todos/{uuid} [put]
func (h *TodosHandler) UpdateTodo(ctx *gin.Context) {
	var request dto.UpdateTodoRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	uuid := ctx.Param("uuid")

	todo, err := h.Usecase.UpdateTodo(ctx, uuid, request)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.NewTodoSuccessResponse(todo, "Update todo successfully")
	ctx.JSON(200, response)
}

// DeleteTodo godoc
// @Summary Delete todo
// @Description Delete todo
// @Tags todos
// @Accept json
// @Produce json
// @security BearerAuth
// @Param uuid path string true "Todo UUID"
// @Success 200 {object} dto.TodoResponse
// @Router /todos/{uuid} [delete]
func (h *TodosHandler) DeleteTodo(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	err := h.Usecase.DeleteTodo(ctx, uuid)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Delete todo successfully",
	})
}
