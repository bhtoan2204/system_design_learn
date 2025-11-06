package repository

import (
	"clean_architect/infrastructure/persistent/model"
	"context"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	ListUsers(ctx context.Context) ([]*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	DeleteUser(ctx context.Context, id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ListUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, user *model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) DeleteUser(ctx context.Context, id string) error {
	return r.db.Delete(&model.User{}, id).Error
}
