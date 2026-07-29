package handler

import (
	"syncova-todo/delivery/dto"
	usecase "syncova-todo/usecase/labels"

	"github.com/gin-gonic/gin"
)

type LabelsHandler struct {
	Usecase usecase.LabelsUseCase
}

func NewLabelsHandler(usecase usecase.LabelsUseCase) *LabelsHandler {
	return &LabelsHandler{
		Usecase: usecase,
	}
}

// GetLabelsByProjectID godoc
// @Summary Get all labels by project ID
// @Description Get all labels by project ID
// @Tags Labels
// @Accept json
// @Produce json
// @security BearerAuth
// @Param project_uuid path string true "Project UUID"
// @Success 200 {object} dto.LabelsListResponse
// @Router /labels/projects/{project_uuid} [get]
func (h *LabelsHandler) GetLabelsByProjectID(ctx *gin.Context) {
	projectUuid := ctx.Param("project_uuid")
	labels, err := h.Usecase.GetLabelsByProjectID(projectUuid)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, dto.NewLabelsListSuccessResponse(labels, "Get labels successfully"))
}

// CreateLabel godoc
// @Summary Create label
// @Description Create label
// @Tags Labels
// @Accept json
// @Produce json
// @security BearerAuth
// @Param request body dto.CreateLabelRequest true "Create label request"
// @Success 201 {object} dto.LabelsResponse
// @Router /labels [post]
func (h *LabelsHandler) CreateLabel(ctx *gin.Context) {
	var label dto.CreateLabelRequest
	if err := ctx.ShouldBindJSON(&label); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	newLabel, err := h.Usecase.CreateLabel(label)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := dto.NewLabelCreateSuccessResponse(newLabel, "Create label successfully")
	ctx.JSON(201, response)
}

// DeleteLabelByUuid godoc
// @Summary Delete label by UUID
// @Description Delete label by UUID
// @Tags Labels
// @Accept json
// @Produce json
// @security BearerAuth
// @Param label_uuid path string true "Label UUID"
// @Success 200 {object} dto.LabelsResponse
// @Router /labels/{label_uuid} [delete]
func (h *LabelsHandler) DeleteLabelByUuid(ctx *gin.Context) {
	labelID := ctx.Param("label_uuid")
	err := h.Usecase.DeleteLabelByUuid(labelID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"success":     true,
		"message":     "Label deleted successfully",
		"status_code": 200,
		"data":        true,
	})
}
