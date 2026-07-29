package http

import (
	"syncova-todo/infrastructure/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.RouterGroup, db *database.PostgresDB) *gin.RouterGroup {

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	var v1 = r.Group("/v1")

	// Setup router
	SetupTodosRouter(v1, db)
	SetupProjectsRouter(v1, db)
	SetupStatusRouter(v1, db)
	SetupLabelsRouter(v1, db)

	return r
}
