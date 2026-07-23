package http

import (
	"syncova-todo/delivery/handler"
	"syncova-todo/infrastructure/database"
	repository "syncova-todo/repository/todos"
	usecase "syncova-todo/usecase/todos"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupTodosRouter(r *gin.RouterGroup, db *database.PostgresDB) *gin.RouterGroup {
	todosRepository := repository.NewTodosRepository(db)
	todosUsecase := usecase.NewTodosUsecase(todosRepository, 5*time.Second)
	var todosRouter = r.Group("/todos")
	todosHandler := handler.NewTodosHandler(todosUsecase)

	todosRouter.GET("/", todosHandler.GetTodos)
	todosRouter.GET("/:uuid", todosHandler.GetTodoByUuid)
	todosRouter.POST("/", todosHandler.CreateTodo)
	todosRouter.PATCH("/:uuid", todosHandler.UpdateTodo)
	todosRouter.DELETE("/:uuid", todosHandler.DeleteTodo)

	return r
}
