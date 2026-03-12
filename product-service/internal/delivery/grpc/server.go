package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"product-service/internal/usecase"
	pb "product-service/pkg/protos/product"
)

type ProductServer struct {
	pb.UnimplementedProductServiceServer
	usecase *usecase.ProductUsecase
}

func NewProductServer(uc *usecase.ProductUsecase) *ProductServer {
	return &ProductServer{usecase: uc}
}

// CreateProduct handles gRPC request to create a product.
func (s *ProductServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product, err := s.usecase.CreateProduct(ctx, req.Name, req.Description, req.Price, req.InitialStock)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}
	return &pb.CreateProductResponse{
		Id:           product.ID.String(),
		Name:         product.Name,
		Description:  product.Description,
		Price:        product.Price,
		InitialStock: product.InitialStock,
	}, nil
}

// GetProduct returns product details without stock.
func (s *ProductServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product id")
	}
	product, err := s.usecase.GetProduct(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}
	if product == nil {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}
	return &pb.GetProductResponse{
		Id:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

// GetProductDetails returns product with current stock (calls Inventory Service).
func (s *ProductServer) GetProductDetails(ctx context.Context, req *pb.GetProductDetailsRequest) (*pb.GetProductDetailsResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product id")
	}
	product, stock, err := s.usecase.GetProductWithStock(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get product details: %v", err)
	}
	if product == nil {
		return nil, status.Errorf(codes.NotFound, "product not found")
	}
	return &pb.GetProductDetailsResponse{
		Id:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Stock:       stock,
	}, nil
}

// UpdateProduct updates a product.
func (s *ProductServer) UpdateProduct(ctx context.Context, req *pb.UpdateProductRequest) (*pb.UpdateProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product id")
	}
	product, err := s.usecase.UpdateProduct(ctx, id, req.Name, req.Description, req.Price)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update product: %v", err)
	}
	return &pb.UpdateProductResponse{
		Id:          product.ID.String(),
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
	}, nil
}

// DeleteProduct deletes a product.
func (s *ProductServer) DeleteProduct(ctx context.Context, req *pb.DeleteProductRequest) (*pb.DeleteProductResponse, error) {
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product id")
	}
	err = s.usecase.DeleteProduct(ctx, id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete product: %v", err)
	}
	return &pb.DeleteProductResponse{Success: true}, nil
}
