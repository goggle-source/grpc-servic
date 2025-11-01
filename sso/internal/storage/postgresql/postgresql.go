package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/goggle-source/grpc-servic/sso/internal/config"
	"github.com/goggle-source/grpc-servic/sso/internal/domain"
	"github.com/goggle-source/grpc-servic/sso/internal/storage"
	"github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(cfg config.Config) (*Storage, error) {
	const op = "postgresql.New"

	conn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Db.User, cfg.Db.Password, cfg.Db.Host, cfg.Db.Port, cfg.Db.NameDB)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil

}

func (s *Storage) SaveUser(ctx context.Context, email string, passwordHash []byte) (int64, error) {
	const op = "postgresql.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	var id int64
	err = stmt.QueryRowContext(ctx, email, passwordHash).Scan(&id)
	if err != nil {
		var psqErr *pq.Error
		if errors.As(err, &psqErr) && psqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (domain.User, error) {
	const op = "postgresql.User"

	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = $1")
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res := stmt.QueryRowContext(ctx, email)

	var user domain.User
	err = res.Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "postgresql.IsAdmin"

	stmt, err := s.db.Prepare("SELECT admin FROM is_admin WHERE id = $1")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	res := stmt.QueryRowContext(ctx, userID)

	var isAdmin bool
	err = res.Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, appID int64) (domain.App, error) {
	const op = "postgresql.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = $1")
	if err != nil {
		return domain.App{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var result domain.App
	res := stmt.QueryRowContext(ctx, appID)
	err = res.Scan(&result.ID, &result.Name, &result.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return domain.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return result, nil
}
