package v1

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"home-library/internal/services/user/dtos"
	"home-library/internal/services/user/usecases"
	customErrors "home-library/pkg/errors"
	"net/http"
)

type handler struct {
	u usecases.UseCase
}

func NewHandler(u usecases.UseCase) *handler {
	return &handler{u: u}
}

func (h *handler) CreateUser(c echo.Context) error {
	var payload dtos.CreateUserRequest
	if err := c.Bind(&payload); err != nil {
		return c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, "Invalid request body", nil))
	}

	if err := payload.Validate(); err != nil {
		validatorErrors := dtos.FromValidatorErrors(err)
		return c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, "Validation failed", validatorErrors))
	}

	userID, err := h.u.CreateUser(context.Background(), payload)
	if err != nil {
		if errors.Is(err, customErrors.ErrUserAlreadyExist) {
			return c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, "User with this email or phone number already exists", nil))
		}
		return c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, "Failed to create user", nil))
	}

	return c.JSON(http.StatusOK, dtos.CreateUserResponse{UserID: userID})
}
