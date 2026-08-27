package domain

import (
	"context"
	"encoding/json"
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
	GetByID(ctx context.Context, id int) (*Product, error)
	Update(ctx context.Context, product *Product) (*Product, error)
}

type ProductService interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error)
	UpdateProduct(ctx context.Context, id int, req *UpdateProductRequest) (*Product, error)
}

type CreateProductRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price"`
}

type UpdateProductRequest struct {
	Name        json.RawMessage `json:"name"`
	Description json.RawMessage `json:"description"`
	SalePrice   json.RawMessage `json:"sale_price"`
	Price       json.RawMessage `json:"price"`
}
