package mocks

import (
	"context"

	"product-api-assessment/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, req *domain.CreateProductRequest) (*domain.Product, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductService) UpdateProduct(ctx context.Context, id int, req *domain.UpdateProductRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}
