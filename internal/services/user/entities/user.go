package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserType string

const (
	UserTypeAdmin UserType = "admin"
	UserTypeUser  UserType = "user"
)

type User struct {
	UserID      uuid.UUID `gorm:"primary_key;type:uuid"`
	FirstName   string    `gorm:"not null;size:50"`
	LastName    string    `gorm:"not null;size:50"`
	Email       string    `gorm:"not null;unique;size:255"`
	PhoneNumber string    `gorm:"not null;unique;size:20"`
	Password    string    `gorm:"not null;size:255"`
	UserType    UserType  `gorm:"not null;type:varchar(5)"`
	IsActive    bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
	DeletedAt   time.Time
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.UserID = uuid.New()
	return nil
}
