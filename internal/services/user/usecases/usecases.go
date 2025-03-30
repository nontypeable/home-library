package usecases

import (
	"context"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"home-library/internal/services/user/dtos"
	"home-library/internal/services/user/entities"
	"home-library/internal/services/user/repository"
	"home-library/pkg/errors"
	"home-library/pkg/jwt"
)

type UseCase interface {
	CreateUser(ctx context.Context, payload dtos.CreateUserRequest) (userID uuid.UUID, err error)
	SignInUser(ctx context.Context, payload dtos.SignInUserRequest) (token string, err error)
}

type useCase struct {
	r   repository.Repository
	jwt jwt.JWTService
}

func NewUseCase(r repository.Repository, jwt jwt.JWTService) UseCase {
	return &useCase{r: r, jwt: jwt}
}

func (u *useCase) CreateUser(ctx context.Context, payload dtos.CreateUserRequest) (userID uuid.UUID, err error) {
	exist, err := u.r.IsUserExist(ctx, payload.Email, payload.PhoneNumber)
	if err != nil {
		return uuid.Nil, err
	}
	if exist {
		return uuid.Nil, errors.ErrUserAlreadyExist
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.Nil, err
	}

	user := entities.NewUser()
	user.FirstName = payload.FirstName
	user.LastName = payload.LastName
	user.Email = payload.Email
	user.PhoneNumber = payload.PhoneNumber
	user.Password = string(hashedPassword)
	user.UserType = entities.UserTypeUser
	user.IsActive = true

	return u.r.CreateUser(ctx, user)
}

func (u *useCase) SignInUser(ctx context.Context, payload dtos.SignInUserRequest) (token string, err error) {
	user, err := u.r.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		return "", errors.ErrInvalidCredentials
	}

	if !user.IsActive {
		return "", errors.ErrUserInactive
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))
	if err != nil {
		return "", errors.ErrInvalidCredentials
	}

	token, err = u.jwt.GenerateToken(jwt.PayloadToken{
		UserID: user.UserID,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
