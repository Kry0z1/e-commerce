package grpcserver

import (
	"context"

	prodcatv1 "github.com/Kry0z1/e-commerce/protos/gen/go/product-catalog"
	"google.golang.org/grpc"
)

type serverAPI struct {
	prodcatv1.CatalogServer
}

func (s *serverAPI) CreateListing(ctx context.Context, req *prodcatv1.CreateListingRequest) (*prodcatv1.CreateListingResponse, error) {
	panic("Unimplemented")
}

func (s *serverAPI) DeleteListing(ctx context.Context, req *prodcatv1.DeleteListingRequest) (*prodcatv1.DeleteListingResponse, error) {
	panic("Unimplemented")
}

func (s *serverAPI) GetListing(ctx context.Context, req *prodcatv1.GetListingRequest) (*prodcatv1.GetListingRequest, error) {
	panic("Unimplemented")
}

func (s *serverAPI) UpdateListing(ctx context.Context, req *prodcatv1.UpdateListingRequest) (*prodcatv1.UpdateListingRequest, error) {
	panic("Unimplemented")
}

func New() *serverAPI {
	return &serverAPI{}
}

func Register(gRPCServer *grpc.Server) {
	prodcatv1.RegisterCatalogServer(gRPCServer, &serverAPI{})
}
