package usecase

import (
	valueobject "clean_architect/domain/value_object"
	"clean_architect/infrastructure/persistent/model"
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

func (u *userImplement) CreateUser(ctx context.Context, input CreateUserInput) (*UserOutput, error) {
	email, err := valueobject.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		Email:    email.Value(),
		Password: input.Password,
	}
	if err := u.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}
	return &UserOutput{
		ID:    user.ID,
		Email: email.Value(),
	}, nil
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
