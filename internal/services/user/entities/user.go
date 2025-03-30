package entities

import (
	"time"

	"github.com/google/uuid"
)

type UserType string

const (
	UserTypeAdmin UserType = "admin"
	UserTypeUser  UserType = "user"
)

type User struct {
	UserID      uuid.UUID  `db:"user_id"`
	FirstName   string     `db:"first_name"`
	LastName    string     `db:"last_name"`
	Email       string     `db:"email"`
	PhoneNumber string     `db:"phone_number"`
	Password    string     `db:"password"`
	UserType    UserType   `db:"user_type"`
	IsActive    bool       `db:"is_active"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at,omitempty"`
}

func NewUser() *User {
	now := time.Now()
	return &User{
		UserID:    uuid.New(),
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}
}
