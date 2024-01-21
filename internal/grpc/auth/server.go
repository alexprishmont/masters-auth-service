package authgrpc

import (
	"auth-sso/internal/services/auth"
	"auth-sso/lib/validation"
	"context"
	"errors"
	authssov1 "github.com/alexprishmont/masters-protos/gen/go/auth-sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type Auth interface {
	Login(ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(ctx context.Context,
		email string,
		password string,
	) (userID string, err error)
	Authorize(ctx context.Context,
		permission string,
		userId string,
	) (isAuthorized bool, err error)
}

type serverAPI struct {
	authssov1.UnimplementedAuthServer
	log  *slog.Logger
	auth Auth
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

type AuthorizeRequest struct {
	Permission string `validate:"required"`
	UserId     string `validate:"required"`
}

func Register(gRPC *grpc.Server, log *slog.Logger, auth Auth) {
	authssov1.RegisterAuthServer(gRPC, &serverAPI{
		log:  log,
		auth: auth,
	})
}

func (s *serverAPI) Login(
	ctx context.Context,
	request *authssov1.LoginRequest,
) (*authssov1.LoginResponse, error) {
	req := LoginRequest{
		Email:    request.GetEmail(),
		Password: request.GetPassword(),
		AppID:    request.GetAppId(),
	}

	if errStr := validation.ValidateStruct(req); errStr != "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", errStr)
	}

	token, err := s.auth.Login(ctx, req.Email, req.Password, int(req.AppID))

	if err != nil {
		if errors.Is(err, auth.ErrorInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authssov1.LoginResponse{
		Token: token,
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

	userID, err := s.auth.RegisterNewUser(ctx, req.Email, req.Password)

	if err != nil {
		if errors.Is(err, auth.ErrorUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) Authorize(
	ctx context.Context,
	request *authssov1.AuthorizeRequest,
) (*authssov1.AuthorizeResponse, error) {
	req := AuthorizeRequest{
		Permission: request.GetPermission(),
		UserId:     request.GetUserId(),
	}

	if errStr := validation.ValidateStruct(req); errStr != "" {
		return nil, status.Errorf(codes.InvalidArgument, "%v", errStr)
	}

	isAuthorized, err := s.auth.Authorize(ctx, req.Permission, req.UserId)

	if err != nil {
		if errors.Is(err, auth.ErrorUserNotAuthorized) {
			return nil, status.Error(codes.Unauthenticated, "Unauthorized action")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &authssov1.AuthorizeResponse{
		Can:        isAuthorized,
		UserId:     req.UserId,
		Permission: req.Permission,
	}, nil
}
