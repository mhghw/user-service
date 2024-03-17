package app

import (
	"github.com/mhghw/user-service/pkg/app/command"
	"github.com/mhghw/user-service/pkg/app/query"
)

type Application struct {
	Commands Command
	Queries  Query
}

type Command struct {
	CreateUser     command.CreateUserHandler
	ChangeUsername command.ChangeUsernameHandler
	DeleteUser     command.DeleteUserHandler
}

type Query struct {
	GetUser query.GetUserHandler
}
