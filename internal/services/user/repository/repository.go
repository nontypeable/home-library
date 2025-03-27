package repository

import (
	"context"
	"home-library/internal/services/user/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(ctx context.Context, user *entities.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	IsUserExist(ctx context.Context, email string) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *entities.User) (uuid.UUID, error) {
	result := r.db.WithContext(ctx).Create(user)
	if result.Error != nil {
		return uuid.Nil, result.Error
	}
	return user.UserID, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *repository) IsUserExist(ctx context.Context, email string) (bool, error) {
	var count int64
	result := r.db.WithContext(ctx).Model(&entities.User{}).Where("email = ?", email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
