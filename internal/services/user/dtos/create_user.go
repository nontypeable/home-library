package dtos

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type CreateUserRequest struct {
	FirstName   string `json:"first_name" validate:"required,min=2,max=50"`
	LastName    string `json:"last_name" validate:"required,min=2,max=50"`
	PhoneNumber string `json:"phone_number" validate:"required,e164"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
}

type CreateUserResponse struct {
	UserID uuid.UUID `json:"user_id"`
}

func (r *CreateUserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
