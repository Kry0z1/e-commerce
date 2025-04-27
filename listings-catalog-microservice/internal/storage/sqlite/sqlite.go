package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"

	"github.com/Kry0z1/e-commerce/listings-catalog-microservice/internal/models"
	"github.com/Kry0z1/e-commerce/listings-catalog-microservice/internal/storage"

	_ "github.com/mattn/go-sqlite3"
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

func (s *Storage) SaveListing(
	ctx context.Context,
	title string,
	description string,
	quantity int64,
	category string,
	closed bool,
	price int64,
	creator int64,
) (int64, error) {
	const op = "storage.sqlite.SaveListing"

	res, err := s.db.ExecContext(ctx, `
		INSERT INTO listings(title, description, quantity, category, closed, price, creator)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, title, description, quantity, category, closed, price, creator)

	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintForeignKey {
			return -1, storage.ErrUserNotFound
		}
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Storage) Listing(ctx context.Context, id int64) (models.Listing, error) {
	const op = "storage.sqlite.Listing"

	var prod models.Listing

	err := s.db.QueryRowContext(ctx, `
		SELECT title, description, quantity, category, closed, price, creator
		FROM listings
		WHERE id = ?
	`, id).Scan(prod.Title, prod.Description, prod.Quantity, prod.Category, prod.Closed, prod.Price, prod.Creator)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return prod, storage.ErrListingNotFound
		}

		return prod, fmt.Errorf("%s: %w", op, err)
	}

	return prod, nil
}

// Nil pointer -> value is unchanged
func (s *Storage) UpdateListing(
	ctx context.Context,
	id int64,
	title *string,
	description *string,
	quantity *int64,
	category *string,
	closed *bool,
	price *int64,
) error {
	const op = "storage.sqlite.UpdateListing"

	res, err := s.db.ExecContext(ctx, `
        UPDATE listings
        SET 
            title = COALESCE(?, title),
            description = COALESCE(?, description),
            quantity = COALESCE(?, quantity),
            category = COALESCE(?, category),
            closed = COALESCE(?, closed),
            price = COALESCE(?, price)
        WHERE id = ?
    `, title, description, quantity, category, closed, price, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrListingNotFound
	}

	return nil
}

func (s *Storage) DeleteListing(ctx context.Context, id int64) error {
	const op = "storage.sqlite.DeleteListing"

	res, err := s.db.ExecContext(ctx, `
        DELETE FROM listings
        WHERE id = ?
    `, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrListingNotFound
	}

	return nil
}
