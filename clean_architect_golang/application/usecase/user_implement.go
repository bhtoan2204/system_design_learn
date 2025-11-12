package usecase

import (
	"clean_architect/infrastructure/persistent/repository"
	"context"
)

type userImplement struct {
	userRepository repository.UserRepository
}

func NewUserImplement(repos repository.Repos) UserUseCase {
	return &userImplement{
		userRepository: repos.UserRepository(),
	}
}

func (u *userImplement) GetUserByID(ctx context.Context, id string) (*UserOutput, error) {
	panic("not implemented")
}

func (u *userImplement) ListUsers(ctx context.Context) ([]*UserOutput, error) {
	panic("not implemented")
}

func (u *userImplement) UpdateUser(ctx context.Context, input UpdateUserInput) (*UserOutput, error) {
	panic("not implemented")
}

func (u *userImplement) DeleteUser(ctx context.Context, id string) error {
	panic("not implemented")
}
