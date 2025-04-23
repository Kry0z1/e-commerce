package service

import (
	"context"
	"log/slog"

	"github.com/Kry0z1/e-commerce/product-catalog-microservice/internal/models"
)

type ProductSaver interface {
	SaveProduct(ctx context.Context) (int64, error)
}

type ProductProvider interface {
	Product(ctx context.Context) (models.Product, error)
}

type Service struct {
	log             *slog.Logger
	productSaver    ProductSaver
	productProvider ProductProvider
}

func New(log *slog.Logger, productSaver ProductSaver, productProvider ProductProvider) *Service {
	return &Service{
		log:             log,
		productSaver:    productSaver,
		productProvider: productProvider,
	}
}
