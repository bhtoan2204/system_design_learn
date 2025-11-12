package handler

import (
	"clean_architect/application/usecase"
	"clean_architect/presentation/http/dto"
	"clean_architect/presentation/http/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase usecase.UserUseCase
}

func NewUserHandler(userUseCase usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := middleware.RequestIDFromCtx(ctx)

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Invalid user ID",
			Message:   "User ID is required",
		})
		return
	}

	user, err := h.userUseCase.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "User not found",
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		RequestID: requestID,
		Data: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

func (h *UserHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := middleware.RequestIDFromCtx(ctx)

	users, err := h.userUseCase.ListUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Failed to list users",
			Message:   err.Error(),
		})
		return
	}

	userResponses := make([]dto.UserResponse, 0, len(users))
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		RequestID: requestID,
		Data:      userResponses,
	})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := middleware.RequestIDFromCtx(ctx)

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Invalid user ID",
			Message:   "User ID is required",
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Invalid request body",
			Message:   err.Error(),
		})
		return
	}

	user, err := h.userUseCase.UpdateUser(ctx, usecase.UpdateUserInput{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Failed to update user",
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		RequestID: requestID,
		Data: dto.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()
	requestID := middleware.RequestIDFromCtx(ctx)

	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Invalid user ID",
			Message:   "User ID is required",
		})
		return
	}

	err := h.userUseCase.DeleteUser(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			RequestID: requestID,
			Error:     "Failed to delete user",
			Message:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.SuccessResponse{
		RequestID: requestID,
		Data:      gin.H{"message": "User deleted successfully"},
	})
}
