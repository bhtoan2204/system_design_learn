package usecase

import "context"

type AuthUseCase interface {
	Register(ctx context.Context, input RegisterInput) (*AuthOutput, error)
	Login(ctx context.Context, input LoginInput) (*AuthOutput, error)
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)
}

type RegisterInput struct {
	Username string
	Email    string
	Password string
}

type LoginInput struct {
	Email    string
	Password string
}

type AuthOutput struct {
	UserID   string
	Username string
	Email    string
	Token    string
}

type TokenClaims struct {
	UserID   string
	Username string
	Email    string
}
