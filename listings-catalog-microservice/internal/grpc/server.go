package grpcserver

import (
	"context"
	"errors"
	"github.com/Kry0z1/e-commerce/listings-catalog-microservice/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	prodcatv1 "github.com/Kry0z1/e-commerce/protos/gen/go/listings-catalog"
	"google.golang.org/grpc"
)

type serverAPI struct {
	prodcatv1.UnimplementedCatalogServer
	srvc service.Service
}

func parseServiceError(err error) error {
	if err != nil {
		if errors.Is(err, service.ErrListingNotFound) || errors.Is(err, service.ErrUserNotFound) {
			return status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, service.ErrNotEnoughPermissions) {
			return status.Error(codes.PermissionDenied, err.Error())
		}
		if errors.Is(err, service.ErrInvalidToken) || errors.Is(err, service.ErrTokenExpired) {
			return status.Error(codes.InvalidArgument, err.Error())
		}

		return status.Error(codes.Internal, "internal error")
	}

	return nil
}

func (s *serverAPI) CreateListing(ctx context.Context, req *prodcatv1.CreateListingRequest) (*prodcatv1.CreateListingResponse, error) {
	title := req.GetTitle()
	if title == "" {
		return nil, status.Error(codes.InvalidArgument, "missing title")
	}

	description := req.GetDescription()
	if description == "" {
		return nil, status.Error(codes.InvalidArgument, "missing description")
	}

	quantity := req.GetQuantity()
	if quantity < 0 {
		return nil, status.Error(codes.InvalidArgument, "quantity cannot be less than 0 dollars")
	}

	category := req.GetCategory()
	if category == "" {
		return nil, status.Error(codes.InvalidArgument, "missing category")
	}

	closed := req.GetClosed()
	price := req.GetPrice()
	if price < 0 {
		return nil, status.Error(codes.InvalidArgument, "price cannot be less than 0 dollars")
	}

	token := req.GetToken()

	id, err := s.srvc.CreateListing(ctx, title, description, quantity, category, closed, price, token)

	return &prodcatv1.CreateListingResponse{Id: id}, parseServiceError(err)
}

func (s *serverAPI) DeleteListing(ctx context.Context, req *prodcatv1.DeleteListingRequest) (*prodcatv1.DeleteListingResponse, error) {
	id := req.GetId()
	token := req.GetToken()

	err := s.srvc.DeleteListing(ctx, id, token)
	if err != nil {
		return &prodcatv1.DeleteListingResponse{Succeeded: false}, parseServiceError(err)
	}

	return &prodcatv1.DeleteListingResponse{Succeeded: true}, nil
}

func (s *serverAPI) GetListing(ctx context.Context, req *prodcatv1.GetListingRequest) (*prodcatv1.GetListingResponse, error) {
	id := req.GetId()

	listing, err := s.srvc.GetListing(ctx, id)

	return &prodcatv1.GetListingResponse{
		Title:       listing.Title,
		Description: listing.Description,
		Quantity:    listing.Quantity,
		Category:    listing.Category,
		Closed:      listing.Closed,
		Price:       listing.Price,
		Creator:     listing.Creator,
	}, parseServiceError(err)
}

func (s *serverAPI) UpdateListing(ctx context.Context, req *prodcatv1.UpdateListingRequest) (*prodcatv1.UpdateListingResponse, error) {
	title := req.GetTitle()
	if title == "" {
		return nil, status.Error(codes.InvalidArgument, "missing title")
	}

	description := req.GetDescription()

	quantity := req.GetQuantity()
	category := req.GetCategory()
	if category == "" {
		return nil, status.Error(codes.InvalidArgument, "missing category")
	}

	closed := req.GetClosed()
	price := req.GetPrice()
	if price < 0 {
		return nil, status.Error(codes.InvalidArgument, "price cannot be less than 0 dollars")
	}

	token := req.GetToken()
	id := req.GetId()

	var descriptionPtr *string
	if description != "" {
		descriptionPtr = &description
	}

	err := s.srvc.UpdateListing(ctx, id, &title, descriptionPtr, &quantity, &category, &closed, &price, token)

	if err != nil {
		return &prodcatv1.UpdateListingResponse{Succeeded: false}, parseServiceError(err)
	}

	return &prodcatv1.UpdateListingResponse{Succeeded: true}, nil
}

func Register(gRPCServer *grpc.Server, srvc service.Service) {
	prodcatv1.RegisterCatalogServer(gRPCServer, &serverAPI{srvc: srvc})
}
