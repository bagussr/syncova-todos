package handler

import (
	"strconv"
	"syncova-todo/delivery/dto"
	domain "syncova-todo/domain/base"
	"syncova-todo/middleware"
	usecase "syncova-todo/usecase/projects"

	"github.com/gin-gonic/gin"
)

type ProjectsHandler struct {
	Usecase usecase.ProjectsUsecase
}

func NewProjectsHandler(usecase usecase.ProjectsUsecase) *ProjectsHandler {
	return &ProjectsHandler{
		Usecase: usecase,
	}
}

// GetProjects godoc
// @Summary Get all projects
// @Description Get all projects
// @Tags Projects
// @Accept json
// @Produce json
// @security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Number of items per page" default(10)
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort query string false "Sort by field" default(desc)
// @Param search query string false "Search by name or description"
// @Success 200 {object} dto.ProjectsListResponse
// @Router /projects [get]
func (h *ProjectsHandler) GetProjects(ctx *gin.Context) {
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

	projects, err := h.Usecase.GetProjects(ctx, request, userID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.NewProjectsListSuccessResponse(projects, "Get projects successfully")

	ctx.JSON(200, response)
}

// GetProjectByUUID godoc
// @Summary Get project by UUID
// @Description Get project by UUID
// @Tags Projects
// @Accept json
// @Produce json
// @security BearerAuth
// @Param uuid path string true "Project UUID"
// @Success 200 {object} dto.ProjectResponse
// @Router /projects/{uuid} [get]
func (h *ProjectsHandler) GetProjectByUUID(ctx *gin.Context) {
	uid := ctx.Param("uuid")

	project, err := h.Usecase.GetProjectByUUID(ctx, uid)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.NewProjectSuccessResponse(project, "Get project successfully")

	ctx.JSON(200, response)
}

// CreateProject godoc
// @Summary Create a new project
// @Description Create a new project
// @Tags Projects
// @Accept json
// @Produce json
// @security BearerAuth
// @Param request body dto.CreateProjectRequest true "Create Project Request"
// @Success 201 {object} dto.ProjectResponse
// @Router /projects [post]
func (h *ProjectsHandler) CreateProject(ctx *gin.Context) {
	var request dto.CreateProjectRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	userId, _ := middleware.GetUserID(ctx)

	project, err := h.Usecase.CreateProject(ctx, request, userId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.NewProjectCreateSuccessResponse(project, "Create project successfully")

	ctx.JSON(201, response)
}

// UpdateProject godoc
// @Summary Update a project
// @Description Update a project
// @Tags Projects
// @Accept json
// @Produce json
// @security BearerAuth
// @Param uuid path string true "Project UUID"
// @Param request body dto.UpdateProjectRequest true "Update Project Request"
// @Success 200 {object} dto.ProjectResponse
// @Router /projects/{uuid} [patch]
func (h *ProjectsHandler) UpdateProject(ctx *gin.Context) {
	uid := ctx.Param("uuid")

	var request dto.UpdateProjectRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	project, err := h.Usecase.UpdateProject(ctx, uid, request)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.NewProjectSuccessResponse(project, "Update project successfully")

	ctx.JSON(200, response)
}

// DeleteProject godoc
// @Summary Delete a project
// @Description Delete a project
// @Tags Projects
// @Accept json
// @Produce json
// @security BearerAuth
// @Param uuid path string true "Project UUID"
// @Success 200 {object} domain.BaseResponse
// @Router /projects/{uuid} [delete]
func (h *ProjectsHandler) DeleteProject(ctx *gin.Context) {
	uid := ctx.Param("uuid")

	err := h.Usecase.DeleteProject(ctx, uid)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, domain.BaseResponse{
		Success:    true,
		StatusCode: 200,
		Message:    "Delete project successfully",
	})
}
