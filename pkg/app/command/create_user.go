package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/mhghw/user-service/decorator"
	apperror "github.com/mhghw/user-service/error"
	"github.com/mhghw/user-service/pkg/domain/user"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

var (
	ErrUsernameAlreadyExists = errors.New("username already exists")
)

type CreateUser struct {
	Username  string
	Firstname string
	Lastname  string
}

type CreateUserHandler decorator.CommandWithResultHandler[CreateUser, string]

type createUserHandler struct {
	userRepo user.Repository
}

func NewCreateUserHandler(userRepo user.Repository, logger *logrus.Entry) CreateUserHandler {
	if userRepo == nil {
		panic("nil user repository")
	}

	return decorator.ApplyCommandWithResultDecorators[CreateUser, string](
		createUserHandler{
			userRepo: userRepo,
		},
		logger,
	)
}

func (h createUserHandler) Handle(ctx context.Context, input CreateUser) (string, error) {
	is, err := h.userRepo.IsUsernameAvailable(ctx, input.Username)
	if err != nil {
		return "", apperror.NewApplicationError(fmt.Errorf("getting availability of username failed: %w", err))
	}
	if !is {
		return "", apperror.NewAlreadyExistsError(ErrUsernameAlreadyExists)
	}

	userID := uuid.New().String()
	user, err := user.New(userID, input.Username, input.Firstname, input.Lastname)
	if err != nil {
		return "", apperror.NewIncorrectInputError(fmt.Errorf("failed to make a new user: %w", err))
	}

	err = h.userRepo.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("error creating user: %w", err)
	}

	return userID, nil
}
