package usecases

import (
	"context"
	"home-library/internal/services/user/dtos"
	"home-library/internal/services/user/entities"
	"home-library/internal/services/user/repository"
	"home-library/pkg/errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UseCase interface {
	CreateUser(ctx context.Context, payload dtos.CreateUserRequest) (userID uuid.UUID, err error)
}

type useCase struct {
	r repository.Repository
}

func NewUseCase(r repository.Repository) UseCase {
	return &useCase{r: r}
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
