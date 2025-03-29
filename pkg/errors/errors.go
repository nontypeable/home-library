package errors

import "errors"

var (
	ErrUserAlreadyExist = errors.New("user with this email or phone number already exist")
)
