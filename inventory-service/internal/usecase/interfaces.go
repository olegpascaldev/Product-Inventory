package usecase

import (
	"context"
	"inventory-service/internal/domain"

	"github.com/google/uuid"
)

type StockRepository interface {
	Create(ctx context.Context, stock *domain.Stock) error
	GetByProductID(ctx context.Context, productID uuid.UUID) (*domain.Stock, error)
	Update(ctx context.Context, stock *domain.Stock) error
	Delete(ctx context.Context, productID uuid.UUID) error
}
