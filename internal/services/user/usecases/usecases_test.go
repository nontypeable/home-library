package usecases

import (
	"context"
	"errors"
	"home-library/internal/services/user/dtos"
	"home-library/internal/services/user/entities"
	customErrors "home-library/pkg/errors"
	"home-library/pkg/jwt"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user *entities.User) (uuid.UUID, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*entities.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.User), args.Error(1)
}

func (m *MockRepository) IsUserExist(ctx context.Context, email string, phoneNumber string) (bool, error) {
	args := m.Called(ctx, email, phoneNumber)
	return args.Bool(0), args.Error(1)
}

type MockJWT struct {
	mock.Mock
}

func (m *MockJWT) GenerateToken(payload jwt.PayloadToken) (string, error) {
	args := m.Called(payload)
	return args.String(0), args.Error(1)
}

func (m *MockJWT) VerifyToken(c echo.Context, token string) error {
	args := m.Called(c, token)
	return args.Error(0)
}

func TestCreateUser(t *testing.T) {
	t.Run("successfully create user", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJWT := new(MockJWT)
		useCase := NewUseCase(mockRepo, mockJWT)
		userID := uuid.New()
		payload := dtos.CreateUserRequest{
			FirstName:   "Evgeny",
			LastName:    "Koveshnikov",
			Email:       "evgeny@example.com",
			PhoneNumber: "+79001234567",
			Password:    "password123",
		}

		mockRepo.On("IsUserExist", mock.Anything, payload.Email, payload.PhoneNumber).
			Return(false, nil)

		mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *entities.User) bool {
			return user.FirstName == payload.FirstName &&
				user.LastName == payload.LastName &&
				user.Email == payload.Email &&
				user.PhoneNumber == payload.PhoneNumber &&
				user.UserType == entities.UserTypeUser &&
				user.IsActive == true
		})).Return(userID, nil)

		id, err := useCase.CreateUser(context.Background(), payload)

		assert.NoError(t, err)
		assert.Equal(t, userID, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user already exists", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJWT := new(MockJWT)
		useCase := NewUseCase(mockRepo, mockJWT)
		payload := dtos.CreateUserRequest{
			FirstName:   "Evgeny",
			LastName:    "Koveshnikov",
			Email:       "evgeny@example.com",
			PhoneNumber: "+79001234567",
			Password:    "password123",
		}

		mockRepo.On("IsUserExist", mock.Anything, payload.Email, payload.PhoneNumber).
			Return(true, nil).Once()

		id, err := useCase.CreateUser(context.Background(), payload)

		assert.Error(t, err)
		assert.Equal(t, customErrors.ErrUserAlreadyExist, err)
		assert.Equal(t, uuid.Nil, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJWT := new(MockJWT)
		useCase := NewUseCase(mockRepo, mockJWT)
		payload := dtos.CreateUserRequest{
			FirstName:   "Evgeny",
			LastName:    "Koveshnikov",
			Email:       "evgeny@example.com",
			PhoneNumber: "+79001234567",
			Password:    "password123",
		}

		expectedErr := errors.New("database error")
		mockRepo.On("IsUserExist", mock.Anything, payload.Email, payload.PhoneNumber).
			Return(false, expectedErr).Once()

		id, err := useCase.CreateUser(context.Background(), payload)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		assert.Equal(t, uuid.Nil, id)
		mockRepo.AssertExpectations(t)
	})

	t.Run("password is properly hashed", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJWT := new(MockJWT)
		useCase := NewUseCase(mockRepo, mockJWT)
		userID := uuid.New()
		payload := dtos.CreateUserRequest{
			FirstName:   "Evgeny",
			LastName:    "Koveshnikov",
			Email:       "evgeny@example.com",
			PhoneNumber: "+79001234567",
			Password:    "password123",
		}

		mockRepo.On("IsUserExist", mock.Anything, payload.Email, payload.PhoneNumber).
			Return(false, nil)

		mockRepo.On("CreateUser", mock.Anything, mock.MatchedBy(func(user *entities.User) bool {
			err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
			return err == nil
		})).Return(userID, nil)

		id, err := useCase.CreateUser(context.Background(), payload)

		assert.NoError(t, err)
		assert.Equal(t, userID, id)
		mockRepo.AssertExpectations(t)
	})
}

func TestSignInUser(t *testing.T) {
	t.Run("successful sign in", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJWT := new(MockJWT)
		useCase := NewUseCase(mockRepo, mockJWT)

		userID := uuid.New()
		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entities.User{
			UserID:   userID,
			Email:    email,
			Password: string(hashedPassword),
			IsActive: true,
		}

		payload := dtos.SignInUserRequest{
			Email:    email,
			Password: password,
		}

		mockRepo.On("GetUserByEmail", context.Background(), email).Return(user, nil)
		mockJWT.On("GenerateToken", jwt.PayloadToken{UserID: userID}).Return("test-token", nil)

		token, err := useCase.SignInUser(context.Background(), payload)

		assert.NoError(t, err)
		assert.Equal(t, "test-token", token)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJWT := new(MockJWT)
		useCase := NewUseCase(mockRepo, mockJWT)

		email := "nonexistent@example.com"
		payload := dtos.SignInUserRequest{
			Email:    email,
			Password: "password123",
		}

		mockRepo.On("GetUserByEmail", context.Background(), email).Return(nil, errors.New("user not found"))

		token, err := useCase.SignInUser(context.Background(), payload)

		assert.Error(t, err)
		assert.Equal(t, customErrors.ErrInvalidCredentials, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertNotCalled(t, "GenerateToken")
	})

	t.Run("inactive account", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJWT := new(MockJWT)
		useCase := NewUseCase(mockRepo, mockJWT)

		userID := uuid.New()
		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entities.User{
			UserID:   userID,
			Email:    email,
			Password: string(hashedPassword),
			IsActive: false,
		}

		payload := dtos.SignInUserRequest{
			Email:    email,
			Password: password,
		}

		mockRepo.On("GetUserByEmail", context.Background(), email).Return(user, nil)

		token, err := useCase.SignInUser(context.Background(), payload)

		assert.Error(t, err)
		assert.Equal(t, customErrors.ErrUserInactive, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertNotCalled(t, "GenerateToken")
	})

	t.Run("invalid password", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJWT := new(MockJWT)
		useCase := NewUseCase(mockRepo, mockJWT)

		userID := uuid.New()
		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entities.User{
			UserID:   userID,
			Email:    email,
			Password: string(hashedPassword),
			IsActive: true,
		}

		payload := dtos.SignInUserRequest{
			Email:    email,
			Password: "wrongpassword",
		}

		mockRepo.On("GetUserByEmail", context.Background(), email).Return(user, nil)

		token, err := useCase.SignInUser(context.Background(), payload)

		assert.Error(t, err)
		assert.Equal(t, customErrors.ErrInvalidCredentials, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertNotCalled(t, "GenerateToken")
	})

	t.Run("jwt generation error", func(t *testing.T) {
		mockRepo := new(MockRepository)
		mockJWT := new(MockJWT)
		useCase := NewUseCase(mockRepo, mockJWT)

		userID := uuid.New()
		email := "test@example.com"
		password := "password123"
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		user := &entities.User{
			UserID:   userID,
			Email:    email,
			Password: string(hashedPassword),
			IsActive: true,
		}

		payload := dtos.SignInUserRequest{
			Email:    email,
			Password: password,
		}

		mockRepo.On("GetUserByEmail", context.Background(), email).Return(user, nil)
		mockJWT.On("GenerateToken", jwt.PayloadToken{UserID: userID}).Return("", errors.New("jwt generation failed"))

		token, err := useCase.SignInUser(context.Background(), payload)

		assert.Error(t, err)
		assert.Empty(t, token)
		mockRepo.AssertExpectations(t)
		mockJWT.AssertExpectations(t)
	})
}
