package mocks

import (
	"context"

	"product-api-assessment/internal/domain"

	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *domain.Product) (*domain.Product, error) {
	args := m.Called(ctx, product)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) GetByID(ctx context.Context, id int) (*domain.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *domain.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}
