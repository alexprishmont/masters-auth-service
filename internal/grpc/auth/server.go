package auth

import (
	"auth-sso/lib/validation"
	"context"
	authssov1 "github.com/alexprishmont/masters-protos/gen/go/auth-sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type serverAPI struct {
	authssov1.UnimplementedAuthServer
	log *slog.Logger
}

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
	AppID    int32  `validate:"required,number,gt=0"`
}

type RegisterRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=6"`
}

func Register(gRPC *grpc.Server, log *slog.Logger) {
	authssov1.RegisterAuthServer(gRPC, &serverAPI{
		log: log,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	request *authssov1.LoginRequest,
) (*authssov1.LoginResponse, error) {
	traceID := ctx.Value("TraceID")

	req := LoginRequest{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
		AppID:    request.GetAppId(),
	}

	if errStr := validation.ValidateStruct(req); errStr != "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", errStr)
	}

	s.log.Info("User successfully signed in.", slog.Int("userID", 32), slog.Any("TraceID", traceID))

	return &authssov1.LoginResponse{
		Token: "Token",
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	request *authssov1.RegisterRequest,
) (*authssov1.RegisterResponse, error) {
	req := RegisterRequest{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
	}

	if errStr := validation.ValidateStruct(req); errStr != "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", errStr)
	}

	return &authssov1.RegisterResponse{
		UserId: 32,
	}, nil
}
