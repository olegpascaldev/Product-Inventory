package usecase

import (
	"context"
	"errors"

	"inventory-service/internal/domain"

	"github.com/google/uuid"
)

type StockUsecase struct {
	repo StockRepository
}

func NewStockUsecase(repo StockRepository) *StockUsecase {
	return &StockUsecase{repo: repo}
}

// CreateStock creates a new stock entry for a product.
func (u *StockUsecase) CreateStock(ctx context.Context, productID uuid.UUID, quantity int32) error {
	// Check if stock already exists
	existing, err := u.repo.GetByProductID(ctx, productID)
	if err != nil {
		return err
	}
	if existing != nil {
		return errors.New("stock already exists for this product")
	}
	stock := domain.NewStock(productID, quantity)
	return u.repo.Create(ctx, stock)
}

// GetStock retrieves stock by product ID.
func (u *StockUsecase) GetStock(ctx context.Context, productID uuid.UUID) (*domain.Stock, error) {
	return u.repo.GetByProductID(ctx, productID)
}

// UpdateStock applies a quantity change (positive or negative) to the stock.
func (u *StockUsecase) UpdateStock(ctx context.Context, productID uuid.UUID, quantityChange int32) (*domain.Stock, error) {
	stock, err := u.repo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, err
	}
	if stock == nil {
		return nil, errors.New("stock not found")
	}
	newQuantity := stock.Quantity + quantityChange
	if newQuantity < 0 {
		return nil, errors.New("insufficient stock")
	}
	stock.Quantity = newQuantity
	if err := u.repo.Update(ctx, stock); err != nil {
		return nil, err
	}
	return stock, nil
}

// DeleteStock removes stock entry.
func (u *StockUsecase) DeleteStock(ctx context.Context, productID uuid.UUID) error {
	return u.repo.Delete(ctx, productID)
}
