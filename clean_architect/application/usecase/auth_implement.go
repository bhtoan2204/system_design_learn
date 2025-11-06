package usecase

import (
	valueobject "clean_architect/domain/value_object"
	"clean_architect/infrastructure/jwt"
	"clean_architect/infrastructure/persistent/model"
	"clean_architect/infrastructure/persistent/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type authImplement struct {
	userRepository repository.UserRepository
	jwtService     jwt.JWTService
}

func NewAuthImplement(repos repository.Repos, jwtService jwt.JWTService) AuthUseCase {
	return &authImplement{
		userRepository: repos.UserRepository(),
		jwtService:     jwtService,
	}
}

func (a *authImplement) Register(ctx context.Context, input RegisterInput) (*AuthOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}

	existingUser, err := a.userRepository.GetUserByEmail(ctx, email.Value())
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	password, err := valueobject.NewPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		ID:        uuid.New().String(),
		Username:  input.Username,
		Email:     email.Value(),
		Password:  password.Value(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := a.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	token, err := a.jwtService.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

func (a *authImplement) Login(ctx context.Context, input LoginInput) (*AuthOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}

	user, err := a.userRepository.GetUserByEmail(ctx, email.Value())
	if err != nil {

		return nil, errors.New("invalid email or password")
	}

	password := valueobject.NewPasswordFromHash(user.Password)
	if !password.Verify(input.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := a.jwtService.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
	}, nil
}

func (a *authImplement) ValidateToken(ctx context.Context, token string) (*TokenClaims, error) {
	claims, err := a.jwtService.ValidateToken(token)
	if err != nil {
		return nil, errors.New("invalid or expired token")
	}

	return &TokenClaims{
		UserID:   claims.UserID,
		Username: claims.Username,
		Email:    claims.Email,
	}, nil
}
