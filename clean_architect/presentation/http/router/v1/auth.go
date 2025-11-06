package router

import (
	"clean_architect/presentation/http/handler"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.RouterGroup, authHandler *handler.AuthHandler) {
	authRouter := router.Group("/auth")
	{
		authRouter.POST("/register", authHandler.Register)
		authRouter.POST("/login", authHandler.Login)
	}
}
