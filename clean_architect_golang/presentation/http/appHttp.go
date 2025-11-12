package http

import (
	"clean_architect/application/usecase"
	"clean_architect/presentation/http/handler"
	"clean_architect/presentation/http/middleware"
	"clean_architect/presentation/http/router/v1"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppHttp interface {
	Routes(ctx context.Context) http.Handler
}

type appHttp struct {
	usecase usecase.Usecase
}

func New(usecase usecase.Usecase) AppHttp {
	return &appHttp{
		usecase: usecase,
	}
}

func (a *appHttp) Routes(ctx context.Context) http.Handler {
	r := gin.New()

	r.Use(middleware.SetRequestID())

	userHandler := handler.NewUserHandler(a.usecase.UserUseCase())
	authHandler := handler.NewAuthHandler(a.usecase.AuthUseCase())

	v1Router := r.Group("/api/v1")
	router.RegisterV1Routes(v1Router, userHandler, authHandler, a.usecase.AuthUseCase())

	return r
}
