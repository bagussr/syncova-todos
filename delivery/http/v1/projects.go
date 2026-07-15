package http

import (
	"syncova-todo/delivery/handler"
	"syncova-todo/infrastructure/database"
	repository "syncova-todo/repository/projects"
	usecase "syncova-todo/usecase/projects"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupProjectsRouter(r *gin.RouterGroup, db *database.PostgresDB) *gin.RouterGroup {
	projectsRepository := repository.NewProjectsRepository(db)
	projectsUsecase := usecase.NewProjectsUsecase(projectsRepository, 30*time.Second)
	var projectsRouter = r.Group("/projects")
	projectsHandler := handler.NewProjectsHandler(projectsUsecase)

	projectsRouter.GET("/", projectsHandler.GetProjects)
	projectsRouter.POST("/", projectsHandler.CreateProject)
	projectsRouter.GET("/:uuid", projectsHandler.GetProjectByUUID)
	projectsRouter.PATCH("/:uuid", projectsHandler.UpdateProject)
	projectsRouter.DELETE("/:uuid", projectsHandler.DeleteProject)

	return r
}
