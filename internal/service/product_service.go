package service

import (
	"context"
	"database/sql"
	"encoding/json"
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

	createdProduct, err := s.repo.Create(ctx, product)
	if err != nil {
		return nil, errors.New(constant.ErrInternalServer)
	}
	return createdProduct, nil
}

func (s *productService) UpdateProduct(ctx context.Context, id int, req *domain.UpdateProductRequest) error {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New(constant.ErrProductNotFound)
		}
		return errors.New(constant.ErrInternalServer)
	}

	if err := mergeUpdateFields(product, req); err != nil {
		return errors.New(constant.ErrInvalidRequest)
	}

	if err := validateProduct(product); err != nil {
		return err
	}

	if err := s.repo.Update(ctx, product); err != nil {
		return errors.New(constant.ErrInternalServer)
	}

	return nil
}

func mergeUpdateFields(product *domain.Product, req *domain.UpdateProductRequest) error {
	if err := applyString(req.Name, &product.Name); err != nil {
		return err
	}
	if err := applyNullableString(req.Description, &product.Description); err != nil {
		return err
	}
	if err := applyNullableFloat(req.SalePrice, &product.SalePrice); err != nil {
		return err
	}
	if err := applyFloat(req.Price, &product.Price); err != nil {
		return err
	}
	return nil
}

func applyString(raw json.RawMessage, target *string) error {
	if raw == nil || string(raw) == "null" {
		return nil
	}
	return json.Unmarshal(raw, target)
}

func applyFloat(raw json.RawMessage, target *float64) error {
	if raw == nil || string(raw) == "null" {
		return nil
	}
	return json.Unmarshal(raw, target)
}

func applyNullableString(raw json.RawMessage, target **string) error {
	if raw == nil {
		return nil
	}
	if string(raw) == "null" {
		*target = nil
		return nil
	}
	var val string
	if err := json.Unmarshal(raw, &val); err != nil {
		return err
	}
	*target = &val
	return nil
}

func applyNullableFloat(raw json.RawMessage, target **float64) error {
	if raw == nil {
		return nil
	}
	if string(raw) == "null" {
		*target = nil
		return nil
	}
	var val float64
	if err := json.Unmarshal(raw, &val); err != nil {
		return err
	}
	*target = &val
	return nil
}

func validateProduct(product *domain.Product) error {
	if product.Name == "" {
		return errors.New(constant.ErrInvalidName)
	}
	if product.Price <= 0 {
		return errors.New(constant.ErrInvalidPrice)
	}
	if product.SalePrice != nil && *product.SalePrice >= product.Price {
		return errors.New(constant.ErrInvalidSalePrice)
	}
	return nil
}
