package errors

import "errors"

var (
	ErrUserAlreadyExist   = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserInactive       = errors.New("user account is inactive")
)
