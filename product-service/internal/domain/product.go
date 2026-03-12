package domain

import "github.com/google/uuid"

// Product represents a product in the domain.
type Product struct {
	ID          uuid.UUID
	Name        string
	Description string
	Price       float64
	//InitialStock используется только при создании товара и не хранится в базе данных.
	InitialStock int32
}

// Функция NewProduct создает новый продукт с cгенерированным UUID.
func CreateProduct(name, description string, price float64, initialStock int32) *Product {
	return &Product{
		ID:           uuid.New(),
		Name:         name,
		Description:  description,
		Price:        price,
		InitialStock: initialStock,
	}
}
