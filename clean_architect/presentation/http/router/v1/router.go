package router

import (
	"clean_architect/presentation/http/handler"

	"github.com/gin-gonic/gin"
)

// RegisterV1Routes registers all v1 API routes
func RegisterV1Routes(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	RegisterUserRoutes(router, userHandler)
}
