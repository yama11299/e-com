package grpc

import (
	context "context"
	"log"
	"net"

	"github.com/yamadev11/e-com/product/internal/app/bl"
	"github.com/yamadev11/e-com/product/internal/app/spec"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server ...
type Server struct {
	UnimplementedProductServer
	product bl.BL
}

// NewGRPCServer ...
func NewGRPCServer(svc bl.BL) *Server {
	return &Server{
		product: svc,
	}
}

// List ...
func (svc *Server) List(ctx context.Context, req *ListRequest) (*ListResponse, error) {
	response := &ListResponse{}
	request := spec.ListRequest{}

	for _, id := range req.Ids {
		request.IDs = append(request.IDs, int(id))
	}

	listResponse, err := svc.product.List(ctx, request)
	if err != nil {
		return nil, err
	}

	for _, product := range listResponse.Products {
		p := getProduct(product)
		response.Products = append(response.Products, &p)
	}

	return response, nil
}

func getProduct(product spec.GetResponse) GetResponse {
	return GetResponse{
		Id:                int32(product.ID),
		Name:              product.Name,
		Price:             int32(product.Price),
		Category:          product.Category,
		AvailableQuantity: int32(product.AvailableQuantity),
	}
}

// UpdateQuantity ...
func (svc *Server) UpdateQuantity(ctx context.Context, req *UpdateQuantityRequest) (*UpdateQuantityResponse, error) {

	updateQuantityRequest := spec.UpdateQuantityRequest{
		ID:       int(req.Id),
		Quantity: int(req.Quantity),
	}

	err := svc.product.UpdateQuantity(ctx, updateQuantityRequest)
	if err != nil {
		return nil, err
	}

	return &UpdateQuantityResponse{}, nil
}

// StartRPCServer ...
func StartRPCServer(svc bl.BL) error {
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("failed to listen !!!")
	}

	productServer := NewGRPCServer(svc)

	server := grpc.NewServer()

	RegisterProductServer(server, productServer)

	reflection.Register(server)

	err = server.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve !!!")
		return err
	}

	return nil
}
