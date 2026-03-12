package usecase

import (
	"context"
	"errors"
	"product-service/internal/domain"

	"github.com/google/uuid"
)

type ProductUsecase struct {
	repo            ProductRepository
	producer        KafkaProducer
	inventoryClient InventoryClient
}

func NewProductUsecase(repo ProductRepository, producer KafkaProducer, inventoryClient InventoryClient) *ProductUsecase {
	return &ProductUsecase{
		repo:            repo,
		producer:        producer,
		inventoryClient: inventoryClient,
	}
}

// Функция CreateProduct создает новый продукт и публикует событие ProductCreated.
func (u *ProductUsecase) CreateProduct(ctx context.Context, name, description string, price float64, initialStock int32) (*domain.Product, error) {
	// Базовая проверка
	if name == "" {
		return nil, errors.New("name is required")
	}
	if price <= 0 {
		return nil, errors.New("price must be positive")
	}
	if initialStock < 0 {
		return nil, errors.New("initial stock cannot be negative")
	}
	product := domain.CreateProduct(name, description, price, initialStock)

	// Сохранение продукта
	if err := u.repo.Create(ctx, product); err != nil {
		return nil, err
	}
	// Публикация события асинхронно (ошибки регистрируются в журнале, но не возвращаются клиенту)
	if err := u.producer.PublishProductCreated(ctx, product.ID, initialStock); err != nil {
		// В реальном приложении используйте логгер
		// Также можно реализовать повторные попытки или очередь недоставленных сообщений
	}
	return product, nil

}

// Функция GetProduct извлекает товар по ID.
func (u *ProductUsecase) GetProduct(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	return u.repo.GetByID(ctx, id)
}

// Метод GetProductWithStock извлекает товар вместе с его текущим наличием на складе из службы учета запасов.
func (u *ProductUsecase) GetProductWithStock(ctx context.Context, id uuid.UUID) (*domain.Product, int32, error) {
	product, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	if product == nil {
		return nil, 0, nil
	}
	stock, err := u.inventoryClient.GetStock(ctx, id)
	if err != nil {
		// Зарегистрировать ошибку, но вернуть товар, которого нет в наличии.
		return product, 0, nil
	}
	return product, stock, nil
}

// Функция UpdateProduct обновляет существующий продукт.
func (u *ProductUsecase) UpdateProduct(ctx context.Context, id uuid.UUID, name, description string, price float64) (*domain.Product, error) {
	product, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("product not found")
	}

	product.Name = name
	product.Description = description
	product.Price = price
	if err := u.repo.Update(ctx, product); err != nil {
		return nil, err
	}
	return product, nil
}

// DeleteProduct удаляет товар.

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id uuid.UUID) error {
	return u.repo.Delete(ctx, id)
}
