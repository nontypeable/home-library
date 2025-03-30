package v1

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
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

func (h *handler) SignInUser(c echo.Context) error {
	var payload dtos.SignInUserRequest
	if err := c.Bind(&payload); err != nil {
		log.Error().Err(err).Msg("failed to bind request body")
		return c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, "Неверный формат запроса", nil))
	}

	if err := payload.Validate(); err != nil {
		validatorErrors := dtos.FromValidatorErrors(err)
		log.Error().Err(err).Interface("validation_errors", validatorErrors).Msg("validation failed")
		return c.JSON(http.StatusBadRequest, dtos.NewErrorResponse(http.StatusBadRequest, "Ошибка валидации", validatorErrors))
	}

	token, err := h.u.SignInUser(context.Background(), payload)
	if err != nil {
		switch {
		case errors.Is(err, customErrors.ErrInvalidCredentials):
			log.Warn().Str("email", payload.Email).Msg("invalid credentials provided")
			return c.JSON(http.StatusUnauthorized, dtos.NewErrorResponse(http.StatusUnauthorized, "Неверный email или пароль", nil))
		case errors.Is(err, customErrors.ErrUserInactive):
			log.Warn().Str("email", payload.Email).Msg("attempt to login with inactive account")
			return c.JSON(http.StatusForbidden, dtos.NewErrorResponse(http.StatusForbidden, "Аккаунт пользователя неактивен", nil))
		default:
			log.Error().Err(err).Str("email", payload.Email).Msg("failed to sign in user")
			return c.JSON(http.StatusInternalServerError, dtos.NewErrorResponse(http.StatusInternalServerError, "Внутренняя ошибка сервера", nil))
		}
	}

	return c.JSON(http.StatusOK, dtos.SignInUserResponse{
		Token: token,
	})
}
