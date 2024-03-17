package user

import "context"

type Repository interface {
	CreateUser(ctx context.Context, u User) error
	UpdateUser(ctx context.Context, userID string, updateFn func(context.Context, User) (User, error)) error
	DeleteUser(ctx context.Context, userID string) error

	GetUser(ctx context.Context, userID string) (User, error)
	IsUsernameAvailable(ctx context.Context, username string) (bool, error)
}
