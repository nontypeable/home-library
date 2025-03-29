package dtos

import (
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Code             int               `json:"code"`
	Message          string            `json:"message"`
	ValidationErrors []ValidationError `json:"validation_errors,omitempty"`
}

type ValidationError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func NewErrorResponse(code int, message string, validationErrors []ValidationError) *ErrorResponse {
	return &ErrorResponse{
		Code:             code,
		Message:          message,
		ValidationErrors: validationErrors,
	}
}

func FromValidatorErrors(err error) []ValidationError {
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	errors := make([]ValidationError, len(validationErrors))
	for i, e := range validationErrors {
		errors[i] = ValidationError{
			Field: e.Field(),
			Tag:   e.Tag(),
			Value: e.Param(),
		}
	}
	return errors
}
