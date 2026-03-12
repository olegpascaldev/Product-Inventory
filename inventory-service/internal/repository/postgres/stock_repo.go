package postgres

import (
	"context"
	"database/sql"
	"errors"

	"inventory-service/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type StockRepository struct {
	db *pgxpool.Pool
}

func NewStockRepository(db *pgxpool.Pool) *StockRepository {
	return &StockRepository{db: db}
}

// Create inserts a new stock record.
func (r *StockRepository) Create(ctx context.Context, stock *domain.Stock) error {
	query := `INSERT INTO stock (product_id, quantity) VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, stock.ProductID, stock.Quantity)
	return err
}

// GetByProductID retrieves stock by product ID.
func (r *StockRepository) GetByProductID(ctx context.Context, productID uuid.UUID) (*domain.Stock, error) {
	query := `SELECT product_id, quantity FROM stock WHERE product_id = $1`
	row := r.db.QueryRow(ctx, query, productID)
	var s domain.Stock
	err := row.Scan(&s.ProductID, &s.Quantity)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &s, nil
}

// Update modifies an existing stock record.
func (r *StockRepository) Update(ctx context.Context, stock *domain.Stock) error {
	query := `UPDATE stock SET quantity = $2 WHERE product_id = $1`
	cmd, err := r.db.Exec(ctx, query, stock.ProductID, stock.Quantity)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// Delete removes a stock record.
func (r *StockRepository) Delete(ctx context.Context, productID uuid.UUID) error {
	query := `DELETE FROM stock WHERE product_id = $1`
	cmd, err := r.db.Exec(ctx, query, productID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}
