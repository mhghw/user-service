package user_test

import (
	"errors"
	"testing"
	"github.com/mhghw/user-service/pkg/domain/user"
)

type NewUserInput struct {
	ID        string
	Username  string
	Firstname string
	Lastname  string
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("New valid user", func(t *testing.T) {
		t.Parallel()

		input := NewUserInput{
			ID:        "1234",
			Username:  "John",
			Firstname: "Doe",
			Lastname:  "johndoe",
		}

		u, _ := user.New(input.ID, input.Username, input.Firstname, input.Lastname)

		assertEqualUser(t, input, u)

	})

	t.Run("User with no username", func(t *testing.T) {
		t.Parallel()

		input := NewUserInput{
			ID:        "1234",
			Username:  "",
			Firstname: "John",
			Lastname:  "Doe",
		}

		u, err := user.New(input.ID, input.Username, input.Firstname, input.Lastname)

		if !errors.Is(err, user.ErrInvalidUsername) {
			t.Errorf("expected error: %v, but got: %v", user.ErrInvalidUsername, err)
		}

		assertEqualUser(t, NewUserInput{}, u)
	})

	t.Run("User with too long username", func(t *testing.T) {
		t.Parallel()

		input := NewUserInput{
			ID:        "1234",
			Username:  "qwerqwerqwerqwerqwerw",
			Firstname: "John",
			Lastname:  "Doe",
		}

		u, err := user.New(input.ID, input.Username, input.Firstname, input.Lastname)

		if !errors.Is(err, user.ErrInvalidUsername) {
			t.Errorf("expected error: %v, but got: %v", user.ErrInvalidUsername, err)
		}

		assertEqualUser(t, NewUserInput{}, u)
	})

	t.Run("User with no firstname", func(t *testing.T) {
		t.Parallel()

		input := NewUserInput{
			ID:        "1234",
			Username:  "johndoe",
			Firstname: "",
			Lastname:  "Doe",
		}

		u, err := user.New(input.ID, input.Username, input.Firstname, input.Lastname)

		if !errors.Is(err, user.ErrInvalidFirstname) {
			t.Errorf("expected error: %v, but got: %v", user.ErrInvalidFirstname, err)
		}

		assertEqualUser(t, NewUserInput{}, u)
	})

	t.Run("User with too long firstname", func(t *testing.T) {
		t.Parallel()

		input := NewUserInput{
			ID:        "1234",
			Username:  "johndoe",
			Firstname: "SmithsonianJonathanQuincyAdamsJohn",
			Lastname:  "Doe",
		}

		u, err := user.New(input.ID, input.Username, input.Firstname, input.Lastname)

		if !errors.Is(err, user.ErrInvalidFirstname) {
			t.Errorf("expected error: %v, but got: %v", user.ErrInvalidFirstname, err)
		}

		assertEqualUser(t, NewUserInput{}, u)
	})

	t.Run("User with no lastname", func(t *testing.T) {
		t.Parallel()

		input := NewUserInput{
			ID:        "1234",
			Username:  "johndoe",
			Firstname: "John",
			Lastname:  "",
		}

		u, err := user.New(input.ID, input.Username, input.Firstname, input.Lastname)

		if !errors.Is(err, user.ErrInvalidLastname) {
			t.Errorf("expected error: %v, but got: %v", user.ErrInvalidLastname, err)
		}

		assertEqualUser(t, NewUserInput{}, u)
	})

	t.Run("User with too long lastname", func(t *testing.T) {
		t.Parallel()

		input := NewUserInput{
			ID:        "1234",
			Username:  "johndoe",
			Firstname: "John",
			Lastname:  "SmithsonianJonathanQuincyAdamsDoe",
		}

		u, err := user.New(input.ID, input.Username, input.Firstname, input.Lastname)

		if !errors.Is(err, user.ErrInvalidLastname) {
			t.Errorf("expected error: %v, but got: %v", user.ErrInvalidLastname, err)
		}

		assertEqualUser(t, NewUserInput{}, u)
	})
}

func assertEqualUser(t *testing.T, input NewUserInput, u user.User) {
	if u.ID != input.ID {
		t.Errorf("id does not match, want: %v, got: %v", input.ID, u.ID)
	}
	if u.Username != input.Username {
		t.Errorf("username does not match, want: %v, got: %v", input.Username, u.Username)
	}
	if u.Firstname != input.Firstname {
		t.Errorf("firstname does not match, want: %v, got: %v", input.Firstname, u.Firstname)
	}
	if u.Lastname != input.Lastname {
		t.Errorf("lastname does not match, want: %v, got: %v", input.Lastname, u.Lastname)
	}
}
