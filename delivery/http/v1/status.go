package http

import (
	"syncova-todo/delivery/handler"
	"syncova-todo/infrastructure/database"
	repository "syncova-todo/repository/status"
	usecase "syncova-todo/usecase/status"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupStatusRouter(r *gin.RouterGroup, db *database.PostgresDB) *gin.RouterGroup {
	statusRepository := repository.NewStatusRepository(db)
	statusUsecase := usecase.NewStatusUsecase(statusRepository, 30*time.Second)
	statusHandler := handler.NewStatusHandler(statusUsecase)
	statusRouter := r.Group("/status")

	statusRouter.GET("/projects/:project_uuid", statusHandler.GetStatusByProjectID)
	statusRouter.GET("/:uuid", statusHandler.GetStatusByUUID)
	statusRouter.POST("/", statusHandler.CreateStatus)
	statusRouter.DELETE("/:uuid", statusHandler.DeleteStatusByUUID)

	return r
}
