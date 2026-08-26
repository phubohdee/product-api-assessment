package service

import (
	"context"
	"errors"

	"product-api-assessment/internal/domain"
	"product-api-assessment/pkg/constant"
)

type productService struct {
	repo domain.ProductRepository
}

func NewProductService(repo domain.ProductRepository) domain.ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error) {
	if req.Name == "" {
		return nil, errors.New(constant.ErrInvalidName)
	}

	if req.Price <= 0 {
		return nil, errors.New(constant.ErrInvalidPrice)
	}

	if req.SalePrice != nil && *req.SalePrice >= req.Price {
		return nil, errors.New(constant.ErrInvalidSalePrice)
	}

	product := &domain.Product{
		Name:        req.Name,
		Description: req.Description,
		SalePrice:   req.SalePrice,
		Price:       req.Price,
	}

	return s.repo.Create(ctx, product)
}
