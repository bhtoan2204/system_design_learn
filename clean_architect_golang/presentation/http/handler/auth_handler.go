package handler

import (
	"clean_architect/application/usecase"
	"clean_architect/presentation/http/dto"
	"clean_architect/presentation/http/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := middleware.RequestIDFromCtx(ctx)

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Invalid request body",
			Message:   err.Error(),
		})
		return
	}

	auth, err := h.authUseCase.Register(ctx, usecase.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		statusCode := http.StatusInternalServerError
		if err.Error() == "user with this email already exists" {
			statusCode = http.StatusConflict
		} else if err.Error() == "invalid email format" || err.Error() == "password must be at least 6 characters" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Registration failed",
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.SuccessResponse{
		RequestID: requestID,
		Data: dto.AuthResponse{
			UserID:   auth.UserID,
			Username: auth.Username,
			Email:    auth.Email,
			Token:    auth.Token,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := middleware.RequestIDFromCtx(ctx)

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Invalid request body",
			Message:   err.Error(),
		})
		return
	}

	auth, err := h.authUseCase.Login(ctx, usecase.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		statusCode := http.StatusUnauthorized
		if err.Error() == "invalid email format" {
			statusCode = http.StatusBadRequest
		}

		c.JSON(statusCode, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Login failed",
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		RequestID: requestID,
		Data: dto.AuthResponse{
			UserID:   auth.UserID,
			Username: auth.Username,
			Email:    auth.Email,
			Token:    auth.Token,
		},
	})
}
