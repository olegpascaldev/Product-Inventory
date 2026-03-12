package grpc

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"inventory-service/internal/usecase"
	pb "inventory-service/pkg/protos/inventory"
)

type InventoryServer struct {
	pb.UnimplementedInventoryServiceServer
	usecase *usecase.StockUsecase
}

func NewInventoryServer(uc *usecase.StockUsecase) *InventoryServer {
	return &InventoryServer{usecase: uc}
}

// GetStock returns the current stock for a product.
func (s *InventoryServer) GetStock(ctx context.Context, req *pb.GetStockRequest) (*pb.GetStockResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product id")
	}
	stock, err := s.usecase.GetStock(ctx, productID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get stock: %v", err)
	}
	if stock == nil {
		return nil, status.Errorf(codes.NotFound, "stock not found")
	}
	return &pb.GetStockResponse{
		ProductId: stock.ProductID.String(),
		Quantity:  stock.Quantity,
	}, nil
}

// UpdateStock modifies the stock quantity.
func (s *InventoryServer) UpdateStock(ctx context.Context, req *pb.UpdateStockRequest) (*pb.UpdateStockResponse, error) {
	productID, err := uuid.Parse(req.ProductId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid product id")
	}
	stock, err := s.usecase.UpdateStock(ctx, productID, req.QuantityChange)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update stock: %v", err)
	}
	return &pb.UpdateStockResponse{
		ProductId:   stock.ProductID.String(),
		NewQuantity: stock.Quantity,
	}, nil
}
