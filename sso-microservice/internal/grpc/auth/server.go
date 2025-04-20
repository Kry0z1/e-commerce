package grpcauth

import (
	"context"
	"errors"

	ssov1 "github.com/Kry0z1/e-commerce/protos/gen/go/sso"
	"github.com/Kry0z1/e-commerce/sso-microservice/internal/services/auth"
	"github.com/Kry0z1/e-commerce/sso-microservice/internal/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string, appID int64) (string, error)
	Register(ctx context.Context, email, password string) (int64, error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
}

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{auth: auth})
}

func (s *serverAPI) RegisterUser(ctx context.Context, req *ssov1.RegisterUserRequest) (*ssov1.RegisterResponse, error) {
	email := req.GetEmail()
	password := req.GetPassword()

	if email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	id, err := s.auth.Register(ctx, email, password)
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.InvalidArgument, "user with such email already exists")
		}
		return nil, status.Error(codes.Internal, "failed to register")
	}

	return &ssov1.RegisterResponse{Id: id}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return nil, status.Error(codes.Internal, "failed to login")
	}

	return &ssov1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if req.UserId == 0 {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Error(codes.Internal, "failed to check admin status")
	}

	return &ssov1.IsAdminResponse{IsAdmin: isAdmin}, nil
}

func New(auth Auth) ssov1.AuthServer {
	return &serverAPI{auth: auth}
}
