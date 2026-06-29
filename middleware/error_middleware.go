package middleware

import (
	"syncova-todo/domain"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if len(ctx.Errors) > 0 {
			err := ctx.Errors.Last().Err
			statusCode, message := domain.ErrorHandler(err)
			ctx.JSON(statusCode, domain.NewErrorResponse(statusCode, message))
		}
	}
}
