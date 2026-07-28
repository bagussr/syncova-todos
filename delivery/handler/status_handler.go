package handler

import (
	"syncova-todo/delivery/dto"
	domain "syncova-todo/domain/base"
	usecase "syncova-todo/usecase/status"

	"github.com/gin-gonic/gin"
)

type StatusHandler struct {
	Usecase usecase.StatusUsecase
}

func NewStatusHandler(usecase usecase.StatusUsecase) *StatusHandler {
	return &StatusHandler{
		Usecase: usecase,
	}
}

// GetStatusByProjectID godoc
// @Summary Get all status by project ID
// @Description Get all status by project ID
// @Tags Status
// @Accept json
// @Produce json
// @security BearerAuth
// @Param project_uuid path string true "Project UUID"
// @Success 200 {object} dto.StatusListResponse
// @Router /status/projects/{project_uuid} [get]
func (h *StatusHandler) GetStatusByProjectID(ctx *gin.Context) {
	uuid, _ := ctx.Params.Get("project_uuid")

	status, err := h.Usecase.GetStatusByProjectID(ctx, uuid)

	if err != nil {
		ctx.JSON(500, dto.NewStatusErrorResponse(err.Error(), 500))
		return
	}

	response := dto.NewStatusListSuccessResponse(status, "Get status successfully")

	ctx.JSON(200, response)

}

// GetStatusByUUID godoc
// @Summary Get status by UUID
// @Description Get status by UUID
// @Tags Status
// @Accept json
// @Produce json
// @security BearerAuth
// @Param uuid path string true "Status UUID"
// @Success 200 {object} dto.StatusResponse
// @Router /status/{uuid} [get]
func (h *StatusHandler) GetStatusByUUID(ctx *gin.Context) {
	uuid, _ := ctx.Params.Get("uuid")

	status, err := h.Usecase.GetStatusByUUID(ctx, uuid)

	if err != nil {
		ctx.JSON(500, dto.NewStatusErrorResponse(err.Error(), 500))
		return
	}

	response := dto.NewStatusSuccessResponse(status, "Get status successfully")

	ctx.JSON(200, response)

}

// CreateStatus godoc
// @Summary Create a new status
// @Description Create a new status
// @Tags Status
// @Accept json
// @Produce json
// @security BearerAuth
// @Param request body dto.CreateStatusRequest true "Create Status Request"
// @Success 201 {object} dto.StatusResponse
// @Router /status [post]
func (h *StatusHandler) CreateStatus(ctx *gin.Context) {
	var request dto.CreateStatusRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(400, dto.NewStatusErrorResponse(err.Error(), 400))
		return
	}

	status, err := h.Usecase.CreateStatus(ctx, request)

	if err != nil {
		ctx.JSON(500, dto.NewStatusErrorResponse(err.Error(), 500))
		return
	}

	response := dto.NewStatusSuccessResponse(status, "Create status successfully")

	ctx.JSON(201, response)

}

// DeleteStatusByUUID godoc
// @Summary Delete status by UUID
// @Description Delete status by UUID
// @Tags Status
// @Accept json
// @Produce json
// @security BearerAuth
// @Param uuid path string true "Status UUID"
// @Success 200 {object} domain.BaseResponse
// @Router /status/{uuid} [delete]
func (h *StatusHandler) DeleteStatusByUUID(ctx *gin.Context) {
	uuid, _ := ctx.Params.Get("uuid")

	err := h.Usecase.DeleteStatusByUUID(ctx, uuid)

	if err != nil {
		ctx.JSON(500, dto.NewStatusErrorResponse(err.Error(), 500))
		return
	}

	ctx.JSON(200, domain.BaseResponse{
		Success:    true,
		StatusCode: 200,
		Message:    "Delete status successfully",
	})
}
