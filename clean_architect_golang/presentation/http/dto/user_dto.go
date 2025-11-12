package dto

import "time"

type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type UpdateUserRequest struct {
	Username string `json:"username" binding:"omitempty,min=3,max=50"`
	Email    string `json:"email" binding:"omitempty,email"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type SuccessResponse struct {
	RequestID string      `json:"request_id"`
	Data      interface{} `json:"data"`
}

type ErrorResponse struct {
	RequestID string `json:"request_id"`
	Error     string `json:"error"`
	Message   string `json:"message"`
}
