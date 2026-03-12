package domain

import "github.com/google/uuid"

// Stock represents the current stock level of a product.
type Stock struct {
	ProductID uuid.UUID
	Quantity  int32
}

// NewStock creates a new stock entry.
func NewStock(productID uuid.UUID, quantity int32) *Stock {
	return &Stock{
		ProductID: productID,
		Quantity:  quantity,
	}
}
