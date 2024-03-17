package service

import (
	"context"
	"fmt"
	"github.com/mhghw/user-service/logs"
	"github.com/mhghw/user-service/pkg/adapters"
	"github.com/mhghw/user-service/pkg/app"
	"github.com/mhghw/user-service/pkg/app/command"
	"github.com/mhghw/user-service/pkg/app/query"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewApplication(ctx context.Context) app.Application {
	db, err := adapters.ConnectToPostgres(ctx, viper.GetString("db.dsn"))
	if err != nil {
		panic(fmt.Errorf("error connecting to db: %w", err))
	}

	logger := logrus.NewEntry(logs.ApplicationLogger())

	userPostgresRepo := adapters.NewUserPostgresRepository(db)

	return app.Application{
		Commands: app.Command{
			CreateUser:     command.NewCreateUserHandler(userPostgresRepo, logger),
			ChangeUsername: command.NewChangeUsernameHandler(userPostgresRepo, logger),
			DeleteUser:     command.NewDeleteUserHandler(userPostgresRepo, logger),
		},
		Queries: app.Query{
			GetUser: query.NewGetUserHandler(userPostgresRepo, logger),
		},
	}
}
