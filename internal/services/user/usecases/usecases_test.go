package usecases

import (
	"context"
	"errors"
	"home-library/internal/services/user/dtos"
	"home-library/internal/services/user/entities"
	customErrors "home-library/pkg/errors"
	"testing"

	"github.com/google/uuid"
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

func TestCreateUser(t *testing.T) {
	t.Run("successfully create user", func(t *testing.T) {
		mockRepo := new(MockRepository)
		useCase := NewUseCase(mockRepo)
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
		useCase := NewUseCase(mockRepo)
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
		useCase := NewUseCase(mockRepo)
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
		useCase := NewUseCase(mockRepo)
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
