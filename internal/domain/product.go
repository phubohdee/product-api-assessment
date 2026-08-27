package domain

import (
	"context"
	"time"
)

type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	SalePrice   *float64  `json:"sale_price"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error)
}
type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"`
}
