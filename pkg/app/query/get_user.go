package query

import (
	"context"
	"fmt"
	"github.com/mhghw/user-service/decorator"
	apperror "github.com/mhghw/user-service/error"
	userPkg "github.com/mhghw/user-service/pkg/domain/user"

	"github.com/sirupsen/logrus"
)

type GetUser struct {
	UserID string
}

type GetUserHandler decorator.QueryHandler[GetUser, userPkg.User]

type getUserHandler struct {
	readModel GetUserReadModel
}

func NewGetUserHandler(readModel GetUserReadModel, logger *logrus.Entry) GetUserHandler {
	if readModel == nil {
		panic("nil get_user readmodel")
	}

	return decorator.ApplyQueryDecorators[GetUser, userPkg.User](
		getUserHandler{
			readModel: readModel,
		},
		logger,
	)
}

type GetUserReadModel interface {
	GetUser(ctx context.Context, userID string) (userPkg.User, error)
}

func (h getUserHandler) Handle(ctx context.Context, input GetUser) (userPkg.User, error) {
	user, err := h.readModel.GetUser(ctx, input.UserID)
	if err != nil {
		return userPkg.User{}, apperror.NewApplicationError(fmt.Errorf("error getting user: %w", err))
	}

	return user, nil
}
