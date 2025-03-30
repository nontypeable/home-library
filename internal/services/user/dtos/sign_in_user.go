package dtos

import (
	"github.com/go-playground/validator/v10"
)

type SignInUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignInUserResponse struct {
	Token string `json:"token"`
}

func (r *SignInUserRequest) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
