package repository

import (
	"context"
	"home-library/internal/services/user/entities"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateUser(ctx context.Context, user *entities.User) (uuid.UUID, error)
	GetUserByEmail(ctx context.Context, email string) (*entities.User, error)
	IsUserExist(ctx context.Context, email string, phoneNumber string) (bool, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(ctx context.Context, user *entities.User) (uuid.UUID, error) {
	query := `
		INSERT INTO users (
			user_id, first_name, last_name, email, phone_number,
			password, user_type, is_active, created_at, updated_at
		) VALUES (
			:user_id, :first_name, :last_name, :email, :phone_number,
			:password, :user_type, :is_active, :created_at, :updated_at
		)
	`

	_, err := r.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return uuid.Nil, err
	}

	return user.UserID, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	query := `
		SELECT * FROM users 
		WHERE email = $1 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) IsUserExist(ctx context.Context, email string, phoneNumber string) (bool, error) {
	var count int
	query := `
		SELECT COUNT(*) FROM users 
		WHERE email = $1 OR phone_number = $2 AND deleted_at IS NULL
	`

	err := r.db.GetContext(ctx, &count, query, email, phoneNumber)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
