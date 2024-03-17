package command

import (
	"context"
	"fmt"
	"github.com/mhghw/user-service/decorator"
	apperror "github.com/mhghw/user-service/error"
	"github.com/mhghw/user-service/pkg/domain/user"

	"github.com/sirupsen/logrus"
)

type ChangeUsername struct {
	UserID      string
	NewUsername string
}

type ChangeUsernameHandler decorator.CommandHandler[ChangeUsername]

type changeUsernameHandler struct {
	userRepo user.Repository
}

func NewChangeUsernameHandler(userRepo user.Repository, logger *logrus.Entry) ChangeUsernameHandler {
	if userRepo == nil {
		panic("nil user repository")
	}

	return decorator.ApplyCommandDecorators[ChangeUsername](
		changeUsernameHandler{
			userRepo: userRepo,
		},
		logger,
	)
}

func (h changeUsernameHandler) Handle(ctx context.Context, input ChangeUsername) error {
	is, err := h.userRepo.IsUsernameAvailable(ctx, input.NewUsername)
	if err != nil {
		return apperror.NewApplicationError(fmt.Errorf("getting availability of username failed: %w", err))
	}
	if !is {
		return apperror.NewAlreadyExistsError(ErrUsernameAlreadyExists)
	}

	err = h.userRepo.UpdateUser(ctx, input.UserID, func(ctx context.Context, u user.User) (user.User, error) {
		u.Username = input.NewUsername
		return u, nil
	})
	if err != nil {
		return apperror.NewApplicationError(fmt.Errorf("error changing user's username: %w", err))
	}

	return nil
}
