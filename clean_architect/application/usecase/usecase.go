package usecase

import (
	"clean_architect/infrastructure/jwt"
	"clean_architect/infrastructure/persistent/repository"
)

type Usecase interface {
	UserUseCase() UserUseCase
	AuthUseCase() AuthUseCase
}

type usecase struct {
	userUseCase UserUseCase
	authUseCase AuthUseCase
}

func NewUsecase(repos repository.Repos, jwtService jwt.JWTService) Usecase {
	return &usecase{
		userUseCase: NewUserImplement(repos),
		authUseCase: NewAuthImplement(repos, jwtService),
	}
}

func (u *usecase) UserUseCase() UserUseCase {
	return u.userUseCase
}

func (u *usecase) AuthUseCase() AuthUseCase {
	return u.authUseCase
}
