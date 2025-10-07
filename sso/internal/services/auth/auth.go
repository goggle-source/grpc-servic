package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/goggle-source/grpc-servic/sso/internal/domain"
	"github.com/goggle-source/grpc-servic/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	SaveUser(
		ctx context.Context,
		email string,
		password []byte,
	) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (domain.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int64) (domain.App, error)
}

type Auth struct {
	log          *slog.Logger
	userSaver    UserStorage
	userProvider UserProvider
	appProvider  AppProvider
	tokenTTL     time.Duration
}

var (
	ErrInvalidCredentials = errors.New("invalid credentails")
)

// New returns new instance of the Auth servic
func New(
	log *slog.Logger,
	userSaver UserStorage,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		appProvider:  appProvider,
		tokenTTL:     tokenTTL,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int64) (token string, err error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("start is login user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		log.Error("field to login user")
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Error("user not found", slog.Any("err", err))

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		log.Error("field to get user", slog.Any("err", err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Error("invalid credentails", slog.Any("err", err))

		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

}

func (a *Auth) Register(ctx context.Context, email string, password string) (userID int64, err error) {
	const op = "auth.Register"

	log := a.log.With(
		slog.String("op", op),
	)

	log.Info("register servic")

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		log.Error("invalid generate hash password", slog.Any("err", err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, email, passwordHash)
	if err != nil {
		log.Error("field to save user", slog.Any("err", err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil

}

func (a *Auth) IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error) {

}
