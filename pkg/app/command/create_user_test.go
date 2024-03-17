package command_test

import (
	"context"
	"fmt"
	"testing"
	apperror "github.com/mhghw/user-service/error"
	"github.com/mhghw/user-service/pkg/app/command"
	userPkg "github.com/mhghw/user-service/pkg/domain/user"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		Name          string
		Input         command.CreateUser
		deps          createUserDependencies
		ShouldFail    bool
		ExpectedError error
	}{
		{
			Name: "create valid user",
			Input: command.CreateUser{
				Username:  "johndoe",
				Firstname: "John",
				Lastname:  "Doe",
			},
			deps:          newCreateUserDependencies(),
			ShouldFail:    false,
			ExpectedError: nil,
		},
		{
			Name: "already exists username",
			Input: command.CreateUser{
				Username:  "johndoe",
				Firstname: "John",
				Lastname:  "Doe",
			},
			deps:          newCreateUserDepsWithSingeUser(),
			ShouldFail:    true,
			ExpectedError: apperror.NewAlreadyExistsError(command.ErrUsernameAlreadyExists), // nil is just for creating a new error type
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			userID, err := tc.deps.handler.Handle(context.Background(), tc.Input)

			if tc.ShouldFail {
				require.ErrorIs(t, err, tc.ExpectedError)
				return
			}

			assertEqualCreatedUser(t, userID, tc.Input, tc.deps)
		})
	}
}

func assertEqualCreatedUser(t *testing.T, userID string, input command.CreateUser, deps createUserDependencies) {
	createdUser, has := deps.userRepo.Users[userID]
	if !has {
		t.Errorf("user %q was not created", userID)
	}

	assert.Equal(t, userID, createdUser.ID)
	assert.Equal(t, input.Username, createdUser.Username)
	assert.Equal(t, input.Firstname, createdUser.Firstname)
	assert.Equal(t, input.Lastname, createdUser.Lastname)
}

type createUserDependencies struct {
	userRepo *repositoryMock
	handler  command.CreateUserHandler
}

func newCreateUserDependencies() createUserDependencies {
	userRepo := &repositoryMock{make(map[string]userPkg.User)}

	return createUserDependencies{
		userRepo: userRepo,
		handler:  command.NewCreateUserHandler(userRepo, logrus.NewEntry(logrus.StandardLogger())),
	}

}

func newCreateUserDepsWithSingeUser() createUserDependencies {
	userRepo := &repositoryMock{
		map[string]userPkg.User{
			"1234": {
				ID:        "1234",
				Username:  "johndoe",
				Firstname: "John",
				Lastname:  "Doe",
			},
		},
	}

	return createUserDependencies{
		userRepo: userRepo,
		handler:  command.NewCreateUserHandler(userRepo, logrus.NewEntry(logrus.StandardLogger())),
	}
}

type repositoryMock struct {
	Users map[string]userPkg.User
}

func (r *repositoryMock) CreateUser(ctx context.Context, u userPkg.User) error {
	r.Users[u.ID] = u
	return nil
}
func (r *repositoryMock) UpdateUser(ctx context.Context, userID string, updateFn func(context.Context, userPkg.User) (userPkg.User, error)) error {
	user, ok := r.Users[userID]
	if !ok {
		return fmt.Errorf("user %q not found", userID)
	}

	updatedUser, err := updateFn(ctx, user)
	if err != nil {
		return err
	}

	r.Users[userID] = updatedUser

	return nil
}
func (r *repositoryMock) DeleteUser(ctx context.Context, userID string) error {
	delete(r.Users, userID)
	return nil
}
func (r *repositoryMock) GetUser(ctx context.Context, userID string) (userPkg.User, error) {
	panic("implement me!")
}
func (r *repositoryMock) IsUsernameAvailable(ctx context.Context, username string) (bool, error) {
	for _, user := range r.Users {
		if username == user.Username {
			return false, nil
		}
	}

	return true, nil
}
