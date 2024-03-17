package ports

import (
	"context"
	"errors"
	"fmt"
	"net"
	apperror "github.com/mhghw/user-service/error"
	"github.com/mhghw/user-service/logs"
	genUser "github.com/mhghw/user-service/pb/gen"
	"github.com/mhghw/user-service/pkg/app"
	"github.com/mhghw/user-service/pkg/app/command"
	"github.com/mhghw/user-service/pkg/app/query"

	"github.com/golang/protobuf/ptypes/empty"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	genUser.UnimplementedUserServiceServer
	app    app.Application
	logger *logrus.Logger
}

func NewGRPCServer(app app.Application) GRPCServer {
	logger := logs.GRPCLogger()

	return GRPCServer{
		app:    app,
		logger: logger,
	}
}

func (s GRPCServer) CreateUser(ctx context.Context, req *genUser.CreateUserRequest) (*genUser.CreateUserResponse, error) {
	input := command.CreateUser{
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
	}

	userID, err := s.app.Commands.CreateUser.Handle(ctx, input)
	if err != nil {
		return nil, handleAppError(err)
	}

	return &genUser.CreateUserResponse{
		UserId: userID,
	}, nil
}

func (s GRPCServer) ChangeUsername(ctx context.Context, req *genUser.ChangeUsernameRequest) (*empty.Empty, error) {
	err := s.app.Commands.ChangeUsername.Handle(ctx, command.ChangeUsername{
		UserID:      req.UserId,
		NewUsername: req.Username,
	})
	if err != nil {
		return nil, handleAppError(err)
	}

	return &empty.Empty{}, nil
}
func (s GRPCServer) DeleteUser(ctx context.Context, req *genUser.DeleteUserRequest) (*empty.Empty, error) {
	err := s.app.Commands.DeleteUser.Handle(ctx, command.DeleteUser{
		UserID: req.UserId,
	})
	if err != nil {
		return nil, handleAppError(err)
	}

	return &empty.Empty{}, nil
}
func (s GRPCServer) GetUser(ctx context.Context, req *genUser.GetUserRequest) (*genUser.User, error) {
	user, err := s.app.Queries.GetUser.Handle(ctx, query.GetUser{UserID: req.UserId})
	if err != nil {
		return nil, handleAppError(err)
	}

	return EncodeUser(user), nil
}

func handleAppError(err error) error {
	var appError *apperror.AppError
	if errors.As(err, &appError) {
		return status.Error(appError.GRPCStatusCode(), appError.Error())
	}

	return status.Error(codes.Unknown, err.Error())
}

func (s GRPCServer) Run(port string) error {
	logrusEntry := logrus.NewEntry(s.logger)

	grpcServer := grpc.NewServer(
		grpc.ChainStreamInterceptor(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
		),
		grpc.ChainStreamInterceptor(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
		),
	)
	genUser.RegisterUserServiceServer(grpcServer, s)

	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Fatal(err)
	}
	s.logger.WithField("grpcEndpoint", addr).Info("Starting: gRPC Listener")

	return grpcServer.Serve(listen)
}
