package router

import (
	"clean_architect/application/usecase"
	"clean_architect/presentation/http/handler"
	"clean_architect/presentation/http/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterV1Routes(router *gin.RouterGroup, userHandler *handler.UserHandler, authHandler *handler.AuthHandler, authUseCase usecase.AuthUseCase) {
	RegisterAuthRoutes(router, authHandler)

	protectedRouter := router.Group("")
	protectedRouter.Use(middleware.AuthMiddleware(authUseCase))
	{
		RegisterUserRoutes(protectedRouter, userHandler)
	}
}
