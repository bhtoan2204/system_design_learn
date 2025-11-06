package usecase

import (
	"context"
	"time"
)

type UserUseCase interface {
	GetUserByID(ctx context.Context, id string) (*UserOutput, error)
	ListUsers(ctx context.Context) ([]*UserOutput, error)
	UpdateUser(ctx context.Context, input UpdateUserInput) (*UserOutput, error)
	DeleteUser(ctx context.Context, id string) error
}

type UpdateUserInput struct {
	ID       string
	Username string
	Email    string
}

type UserOutput struct {
	ID        string
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
