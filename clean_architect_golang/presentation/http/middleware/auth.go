package middleware

import (
	"clean_architect/application/usecase"
	"clean_architect/constant"
	"clean_architect/presentation/http/dto"
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authUseCase usecase.AuthUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		requestID := RequestIDFromCtx(ctx)

		// Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				RequestID: requestID,
				Error:     "Unauthorized",
				Message:   "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				RequestID: requestID,
				Error:     "Unauthorized",
				Message:   "Invalid authorization header format. Expected: Bearer <token>",
			})
			c.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := authUseCase.ValidateToken(ctx, token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				RequestID: requestID,
				Error:     "Unauthorized",
				Message:   err.Error(),
			})
			c.Abort()
			return
		}

		// Set user info in context
		ctx = withUserID(ctx, claims.UserID)
		ctx = withUsername(ctx, claims.Username)
		ctx = withEmail(ctx, claims.Email)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func UserIDFromCtx(ctx context.Context) string {
	v := ctx.Value(constant.CtxKeyUserID)
	if v == nil {
		return ""
	}
	if val, ok := v.(string); ok {
		return val
	}
	return ""
}

func UsernameFromCtx(ctx context.Context) string {
	v := ctx.Value(constant.CtxKeyUsername)
	if v == nil {
		return ""
	}
	if val, ok := v.(string); ok {
		return val
	}
	return ""
}

func EmailFromCtx(ctx context.Context) string {
	v := ctx.Value(constant.CtxKeyEmail)
	if v == nil {
		return ""
	}
	if val, ok := v.(string); ok {
		return val
	}
	return ""
}

func withUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, constant.CtxKeyUserID, userID)
}

func withUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, constant.CtxKeyUsername, username)
}

func withEmail(ctx context.Context, email string) context.Context {
	return context.WithValue(ctx, constant.CtxKeyEmail, email)
}
