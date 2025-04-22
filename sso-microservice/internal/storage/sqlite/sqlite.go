package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"

	"github.com/Kry0z1/e-commerce/sso-microservice/internal/domain/models"
	"github.com/Kry0z1/e-commerce/sso-microservice/internal/storage"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

func (s *Storage) SaveUser(ctx context.Context, email string, hashedPassword []byte) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	res, err := s.db.ExecContext(ctx, `
		INSERT INTO users(email, pass_hash) VALUES(?, ?)
	`, email, hashedPassword)

	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return -1, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	var user models.User

	err := s.db.QueryRowContext(ctx, `
		SELECT id, email, pass_hash
		FROM users
		WHERE email == ?
	`, email).Scan(&user.ID, &user.Email, &user.HashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, id int64) (bool, error) {
	const op = "storage.sqlite.IsAdmin"

	var isAdmin bool

	err := s.db.QueryRowContext(ctx, `
		SELECT is_admin
		FROM users
		WHERE id == ?
	`, id).Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return isAdmin, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return isAdmin, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) App(ctx context.Context, id int64) (models.App, error) {

	const op = "storage.sqlite.App"

	var app models.App

	err := s.db.QueryRowContext(ctx, `
		SELECT id, name, secret
		FROM apps
		WHERE id == ?
	`, id).Scan(&app.ID, &app.Name, &app.SecretKey)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return app, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return app, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}
