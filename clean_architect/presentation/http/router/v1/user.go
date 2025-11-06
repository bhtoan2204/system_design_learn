package router

import (
	"clean_architect/presentation/http/handler"

	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	userRouter := router.Group("/users")
	{
		userRouter.POST("", userHandler.CreateUser)
		userRouter.GET("", userHandler.ListUsers)
		userRouter.GET("/:id", userHandler.GetUserByID)
		userRouter.PUT("/:id", userHandler.UpdateUser)
		userRouter.DELETE("/:id", userHandler.DeleteUser)
	}
}
