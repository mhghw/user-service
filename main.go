package main

import (
	"context"
	"github.com/mhghw/user-service/config"
	"github.com/mhghw/user-service/pkg/ports"
	"github.com/mhghw/user-service/pkg/service"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	config.Init()

	ctx := context.Background()

	app := service.NewApplication(ctx)

	grpcServer := ports.NewGRPCServer(app)

	logrus.Fatal(grpcServer.Run(viper.GetString("bootstrap.port")))
}
