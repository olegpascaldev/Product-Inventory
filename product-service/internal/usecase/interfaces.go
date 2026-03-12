package usecase

import (
	"context"
	"product-service/internal/domain"

	"github.com/google/uuid"
)

// ProductRepository определяет методы репозитория, необходимые для конкретного варианта использования.
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// KafkaProducer определяет методы для публикации событий предметной области.
type KafkaProducer interface {
	PublishProductCreated(ctx context.Context, productID uuid.UUID, initialStock int32) error
}

// InventoryClient определяет методы для вызова службы инвентаризации через gRPC.
type InventoryClient interface {
	GetStock(ctx context.Context, productID uuid.UUID) (int32, error)
}
