package command

import (
	"context"
	"fmt"
	"github.com/mhghw/user-service/decorator"
	apperror "github.com/mhghw/user-service/error"
	"github.com/mhghw/user-service/pkg/domain/user"

	"github.com/sirupsen/logrus"
)

type DeleteUser struct {
	UserID string
}

type DeleteUserHandler decorator.CommandHandler[DeleteUser]

type deleteUserHandler struct {
	userRepo user.Repository
}

func NewDeleteUserHandler(userRepo user.Repository, logger *logrus.Entry) DeleteUserHandler {
	if userRepo == nil {
		panic("nil user repository")
	}

	return decorator.ApplyCommandDecorators[DeleteUser](
		deleteUserHandler{userRepo: userRepo}, logger,
	)
}

func (h deleteUserHandler) Handle(ctx context.Context, input DeleteUser) error {
	err := h.userRepo.DeleteUser(ctx, input.UserID)
	if err != nil {
		return apperror.NewApplicationError(fmt.Errorf("error deleting user: %w", err))
	}

	return nil
}
