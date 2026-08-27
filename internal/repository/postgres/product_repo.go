package postgres

import (
	"context"
	"database/sql"

	"product-api-assessment/internal/domain"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) domain.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	query := `
		INSERT INTO products (name, description, sale_price, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, description, sale_price, price, created_at, updated_at
	`

	created := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.SalePrice,
		product.Price,
	).Scan(
		&created.ID,
		&created.Name,
		&created.Description,
		&created.SalePrice,
		&created.Price,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return created, nil
}
