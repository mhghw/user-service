package user

import "errors"

var (
	ErrInvalidUsername  = errors.New("invalid username, length must be greater than 4 and less than 20")
	ErrInvalidFirstname = errors.New("invalid firstname, length must be greater than 2 and less than 30")
	ErrInvalidLastname  = errors.New("invalid lastname, length must be greater than 2 and less than 30")
)

type User struct {
	ID        string
	Username  string
	Firstname string
	Lastname  string
}

func New(id, username, firstname, lastname string) (User, error) {
	if username == "" || len(username) > 20 {
		return User{}, ErrInvalidUsername
	}

	if firstname == "" || len(firstname) > 30 {
		return User{}, ErrInvalidFirstname
	}

	if lastname == "" || len(lastname) > 30 {
		return User{}, ErrInvalidLastname
	}

	return User{
		ID:        id,
		Username:  username,
		Firstname: firstname,
		Lastname:  lastname,
	}, nil
}
