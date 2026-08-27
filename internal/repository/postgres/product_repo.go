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

func (r *productRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	query := `SELECT id, name, description, sale_price, price, created_at, updated_at FROM products WHERE id = $1`

	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.SalePrice,
		&product.Price,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	query := `
		UPDATE products
		SET name = $1, description = $2, sale_price = $3, price = $4, updated_at = NOW()
		WHERE id = $5
		RETURNING id, name, description, sale_price, price, created_at, updated_at
	`

	updated := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Description,
		product.SalePrice,
		product.Price,
		product.ID,
	).Scan(
		&updated.ID,
		&updated.Name,
		&updated.Description,
		&updated.SalePrice,
		&updated.Price,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return updated, nil
}
