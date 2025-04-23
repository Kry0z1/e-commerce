package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Kry0z1/e-commerce/product-catalog-microservice/internal/models"

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

func (s *Storage) SaveProduct(ctx context.Context) (int64, error) {
	panic("Unimplemented")
}

func (s *Storage) Product(ctx context.Context) (models.Product, error) {
	panic("Unimplemented")
}
