package auth

import (
	"context"
	"strings"

	ssov1 "github.com/goggle-source/grpc-servic/protos/gen/go/sso"
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

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))

	if err != nil {
		//TODO: add proverok
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
		//TODO: ...
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
		//TODO: ...
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
	if req.GetPassword() == "" || len(req.GetPassword()) > 5 {
		return status.Error(codes.InvalidArgument, "password is required")
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
	if req.GetPassword() == "" || len(req.GetPassword()) > 5 {
		return status.Error(codes.InvalidArgument, "password is required")
	}

	return nil
}

func ValidateIsAdmin(req *ssov1.IsAdminRequest) error {
	if req.GetUserId() == emptyID {
		return status.Error(codes.InvalidArgument, "user_id is requred")
	}

	return nil
}
