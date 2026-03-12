package postgres

import (
	"context"
	"database/sql"
	"errors"
	"product-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ProductRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) *ProductRepository {
	return &ProductRepository{
		db: db,
	}
}

// Функция Create вставляет новый товар в базу данных.
func (r *ProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (id, name, description, price) VALUES ($1, $2, $3, $4)`
	_, err := r.db.Exec(ctx, query, product.ID, product.Name, product.Description, product.Price)
	return err
}

// Функция GetByID извлекает товар по его идентификатору.
func (r *ProductRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `SELECT id, name, description, price FROM products WHERE id = $1`
	row := r.db.QueryRow(ctx, query, id)
	var p domain.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// Функция обновления обновляет существующий продукт.
func (r *ProductRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `UPDATE products SET name = $2, description = $3, price = $4 WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, product.ID, product.Name, product.Description, product.Price)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *ProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}
