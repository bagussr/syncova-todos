package http

import (
	"syncova-todo/delivery/handler"
	"syncova-todo/infrastructure/database"
	repository "syncova-todo/repository/labels"
	usecase "syncova-todo/usecase/labels"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupLabelsRouter(r *gin.RouterGroup, db *database.PostgresDB) *gin.RouterGroup {
	labelsRepository := repository.NewLabelsRepository(db)
	labelsUsecase := usecase.NewLabelsUsecase(labelsRepository, 5*time.Second)
	var labelsRouter = r.Group("/labels")
	labelsHandler := handler.NewLabelsHandler(labelsUsecase)

	labelsRouter.GET("/projects/:project_uuid", labelsHandler.GetLabelsByProjectID)
	labelsRouter.POST("/", labelsHandler.CreateLabel)
	labelsRouter.DELETE("/:uuid", labelsHandler.DeleteLabelByUuid)

	return r
}
