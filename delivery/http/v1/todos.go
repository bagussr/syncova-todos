package http

import (
	"syncova-todo/delivery/handler"
	repository "syncova-todo/repository/todos"
	usecase "syncova-todo/usecase/todos"
	"time"

	"github.com/gin-gonic/gin"
)

func SetupTodosRouter(r *gin.RouterGroup) *gin.RouterGroup {
	todosRepository := repository.NewTodosRepository()
	todosUsecase := usecase.NewTodosUsecase(todosRepository, 5*time.Second)
	var todosRouter = r.Group("/todos")
	todosHandler := handler.NewTodosHandler(todosUsecase)

	todosRouter.GET("/", todosHandler.GetTodos)

	return r
}
