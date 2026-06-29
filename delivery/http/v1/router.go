package http

import "github.com/gin-gonic/gin"

func SetupRouter(r *gin.RouterGroup) *gin.RouterGroup {

	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello, World!",
		})
	})

	var v1 = r.Group("/v1")

	// Setup router
	SetupTodosRouter(v1)

	return r
}
