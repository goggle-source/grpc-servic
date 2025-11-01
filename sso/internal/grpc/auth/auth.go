package Grpcauth

import (
	"context"
	"errors"
	"strings"

	ssov1 "github.com/goggle-source/grpc-servic/protos/gen/go/sso"
	"github.com/goggle-source/grpc-servic/sso/internal/services/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	emptyID = 0
)

type ServicAuth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int64,
	) (token string, err error)

	Register(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)

	IsAdmin(
		ctx context.Context,
		userID int64,
	) (isAdmin bool, err error)
}

type ServerAPI struct {
	ssov1.UnimplementedAuthServer
	auth ServicAuth
}

func Register(gRPC *grpc.Server, auth ServicAuth) {
	ssov1.RegisterAuthServer(gRPC, &ServerAPI{auth: auth})
}

func (s *ServerAPI) Login(ctx context.Context, req *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {

	if err := ValidateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int64(req.GetAppId()))

	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.Internal, "error credentails")
		}
		if errors.Is(err, auth.ErrAppNotFound) {
			return nil, status.Error(codes.NotFound, "app is not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {
	if err := ValidateRegister(req); err != nil {
		return nil, err
	}

	userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, auth.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user alredy exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *ServerAPI) IsAdmin(ctx context.Context, req *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	if err := ValidateIsAdmin(req); err != nil {
		return nil, err
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, auth.ErrAppNotFound) {
			return nil, status.Error(codes.NotFound, "app is not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &ssov1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil
}

func ValidateLogin(req *ssov1.LoginRequest) error {
	if req.GetEmail() == "" || !strings.Contains(req.GetEmail(), "@") {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if len(req.GetPassword()) < 7 {
		return status.Error(codes.InvalidArgument, "password must be at least 10 characters")
	}

	if req.GetAppId() == emptyID {
		return status.Error(codes.InvalidArgument, "app_id is required")
	}

	return nil
}

func ValidateRegister(req *ssov1.RegisterRequest) error {
	if req.GetEmail() == "" || !strings.Contains(req.GetEmail(), "@") {
		return status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	if len(req.GetPassword()) < 7 {
		return status.Error(codes.InvalidArgument, "password must be at least 10 characters")
	}

	return nil
}

func ValidateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyID {
		return status.Error(codes.InvalidArgument, "user_id is requred")
	}

	return nil
}
