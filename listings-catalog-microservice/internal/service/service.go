package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Kry0z1/e-commerce/listings-catalog-microservice/internal/jwt"
	"github.com/Kry0z1/e-commerce/listings-catalog-microservice/internal/models"
	"github.com/Kry0z1/e-commerce/listings-catalog-microservice/internal/storage"
)

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrListingNotFound      = errors.New("listing not found")
	ErrNotEnoughPermissions = errors.New("user is not authorized for this action")
	ErrTokenExpired         = errors.New("token is expired")
	ErrInvalidToken         = errors.New("token is invalid")
)

type ListingSaver interface {
	SaveListing(
		ctx context.Context,
		title string,
		description string,
		quantity int64,
		category string,
		closed bool,
		price int64,
		creator int64,
	) (int64, error)

	// Nil pointer -> value is unchanged
	UpdateListing(
		ctx context.Context,
		id int64,
		title *string,
		description *string,
		quantity *int64,
		category *string,
		closed *bool,
		price *int64,
	) error

	DeleteListing(ctx context.Context, id int64) error
}

type ListingProvider interface {
	Listing(ctx context.Context, id int64) (models.Listing, error)
}

type Service struct {
	log             *slog.Logger
	productSaver    ListingSaver
	productProvider ListingProvider
}

func New(log *slog.Logger, productSaver ListingSaver, productProvider ListingProvider) *Service {
	return &Service{
		log:             log,
		productSaver:    productSaver,
		productProvider: productProvider,
	}
}

func (s *Service) CreateListing(
	ctx context.Context,
	title string,
	description string,
	quantity int64,
	category string,
	closed bool,
	price int64,
	token string,
) (int64, error) {
	const op = "service.CreateListing"

	log := s.log.With(slog.String("op", op))

	log.Info("started listing creation")

	tokenData, err := jwt.ParseToken(token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Info("token expired")
			return -1, ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenInvalid) {
			log.Info("token invalid: %s", err.Error())
			return -1, ErrInvalidToken
		}
		log.Error("failed to parse token: %s", err.Error())
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.productSaver.SaveListing(ctx, title, description, quantity, category, closed, price, tokenData.ID)

	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			log.Info("user not found")
			return -1, ErrUserNotFound
		}
		log.Error("failed to find user: %s", err.Error())
		return -1, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("creation succeeded")
	return id, nil
}

func (s *Service) DeleteListing(ctx context.Context, id int64, token string) error {
	const op = "service.DeleteListing"

	log := s.log.With(slog.String("op", op))

	log.Info("started listing deletion")

	tokenData, err := jwt.ParseToken(token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Info("token expired")
			return ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenInvalid) {
			log.Info("token invalid: %s", err.Error())
			return ErrInvalidToken
		}
		log.Error("failed to parse token: %s", err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}

	listing, err := s.productProvider.Listing(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrListingNotFound) {
			log.Info("listing not found on get")
			return ErrListingNotFound
		}
		log.Error("internal error: %s", err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}

	if listing.Creator != tokenData.ID {
		log.Info("wrong user")
		return ErrNotEnoughPermissions
	}

	if err := s.productSaver.DeleteListing(ctx, id); err != nil {
		if errors.Is(err, storage.ErrListingNotFound) {
			log.Info("listing not found on delete")
			return ErrListingNotFound
		}
		log.Error("internal error: %s", err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("deletion succeeded")
	return nil
}

func (s *Service) GetListing(ctx context.Context, id int64) (models.Listing, error) {
	const op = "service.GetListing"

	log := s.log.With(slog.String("op", op))

	log.Info("started listing getting")

	listing, err := s.productProvider.Listing(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrListingNotFound) {
			log.Info("listing not found on get")
			return listing, ErrListingNotFound
		}
		log.Error("internal error: %s", err.Error())
		return listing, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("getting succeeded")
	return listing, nil
}

// Nil pointer -> value is unchanged
func (s *Service) UpdateListing(
	ctx context.Context,
	id int64,
	title *string,
	description *string,
	quantity *int64,
	category *string,
	closed *bool,
	price *int64,
	token string,
) error {
	const op = "service.UpdateListing"

	log := s.log.With(slog.String("op", op))

	log.Info("started listing updating")

	tokenData, err := jwt.ParseToken(token)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.Info("token expired")
			return ErrTokenExpired
		}
		if errors.Is(err, jwt.ErrTokenInvalid) {
			log.Info("token invalid: %s", err.Error())
			return ErrInvalidToken
		}
		log.Error("failed to parse token: %s", err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}

	listing, err := s.productProvider.Listing(ctx, id)
	if err != nil {
		if errors.Is(err, storage.ErrListingNotFound) {
			log.Info("listing not found on get")
			return ErrListingNotFound
		}
		log.Error("internal error: %s", err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}

	if listing.Creator != tokenData.ID {
		log.Info("wrong user")
		return ErrNotEnoughPermissions
	}

	if err := s.productSaver.UpdateListing(ctx, id, title, description, quantity, category, closed, price); err != nil {
		if errors.Is(err, storage.ErrListingNotFound) {
			log.Info("listing not found on delete")
			return ErrListingNotFound
		}
		log.Error("internal error: %s", err.Error())
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("update succeeded")
	return nil
}
