package v1

import (
	"context"
	"encoding/json"
	"errors"
	"home-library/internal/services/user/dtos"
	customErrors "home-library/pkg/errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUseCase struct {
	mock.Mock
}

func (m *MockUseCase) CreateUser(ctx context.Context, payload dtos.CreateUserRequest) (uuid.UUID, error) {
	args := m.Called(ctx, payload)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	e := echo.New()

	t.Run("successfully create user", func(t *testing.T) {
		mockUseCase := new(MockUseCase)
		handler := NewHandler(mockUseCase)
		userID := uuid.New()

		payload := dtos.CreateUserRequest{
			FirstName:   "Evgeny",
			LastName:    "Koveshnikov",
			Email:       "evgeny@example.com",
			PhoneNumber: "+79001234567",
			Password:    "password123",
		}

		jsonPayload, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUseCase.On("CreateUser", context.Background(), payload).Return(userID, nil)

		err := handler.CreateUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, userID.String(), response["id"])

		mockUseCase.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		mockUseCase := new(MockUseCase)
		handler := NewHandler(mockUseCase)

		payload := dtos.CreateUserRequest{
			FirstName:   "E",
			LastName:    "Koveshnikov",
			Email:       "invalid-email",
			PhoneNumber: "123",
			Password:    "123",
		}

		jsonPayload, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response dtos.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Validation failed", response.Message)
		assert.NotEmpty(t, response.ValidationErrors)

		mockUseCase.AssertNotCalled(t, "CreateUser")
	})

	t.Run("user already exists", func(t *testing.T) {
		mockUseCase := new(MockUseCase)
		handler := NewHandler(mockUseCase)

		payload := dtos.CreateUserRequest{
			FirstName:   "Evgeny",
			LastName:    "Koveshnikov",
			Email:       "evgeny@example.com",
			PhoneNumber: "+79001234567",
			Password:    "password123",
		}

		jsonPayload, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUseCase.On("CreateUser", context.Background(), payload).Return(uuid.Nil, customErrors.ErrUserAlreadyExist)

		err := handler.CreateUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response dtos.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "User with this email or phone number already exists", response.Message)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		mockUseCase := new(MockUseCase)
		handler := NewHandler(mockUseCase)

		payload := dtos.CreateUserRequest{
			FirstName:   "Evgeny",
			LastName:    "Koveshnikov",
			Email:       "evgeny@example.com",
			PhoneNumber: "+79001234567",
			Password:    "password123",
		}

		jsonPayload, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader(string(jsonPayload)))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mockUseCase.On("CreateUser", context.Background(), payload).Return(uuid.Nil, errors.New("database error"))

		err := handler.CreateUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response dtos.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Failed to create user", response.Message)

		mockUseCase.AssertExpectations(t)
	})

	t.Run("invalid json", func(t *testing.T) {
		mockUseCase := new(MockUseCase)
		handler := NewHandler(mockUseCase)

		req := httptest.NewRequest(http.MethodPost, "/sign-up", strings.NewReader("invalid json"))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.CreateUser(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response dtos.ErrorResponse
		err = json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request body", response.Message)

		mockUseCase.AssertNotCalled(t, "CreateUser")
	})
}
