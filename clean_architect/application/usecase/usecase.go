package usecase

type Usecase interface {
	UserUseCase() UserUseCase
}

type usecase struct {
	userUseCase UserUseCase
}

func NewUsecase() Usecase {
	return &usecase{}
}

func (u *usecase) UserUseCase() UserUseCase {
	return u.userUseCase
}
