package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"product-api-assessment/internal/domain"
	"product-api-assessment/mocks"
	"product-api-assessment/pkg/constant"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)

	desc := "Test description"
	req := &domain.CreateProductRequest{
		Name:        "Test Product",
		Description: &desc,
		Price:       100.00,
	}

	expectedProduct := &domain.Product{
		ID:          1,
		Name:        "Test Product",
		Description: &desc,
		Price:       100.00,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(expectedProduct, nil)

	result, err := svc.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedProduct.Name, result.Name)
	assert.Equal(t, expectedProduct.Price, result.Price)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_SuccessWithSalePrice(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)

	salePrice := 80.00
	req := &domain.CreateProductRequest{
		Name:      "Sale Product",
		SalePrice: &salePrice,
		Price:     100.00,
	}

	expectedProduct := &domain.Product{
		ID:        1,
		Name:      "Sale Product",
		SalePrice: &salePrice,
		Price:     100.00,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*domain.Product")).Return(expectedProduct, nil)

	result, err := svc.CreateProduct(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, &salePrice, result.SalePrice)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_EmptyName(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)

	req := &domain.CreateProductRequest{
		Name:  "",
		Price: 100.00,
	}

	result, err := svc.CreateProduct(context.Background(), req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, constant.ErrInvalidName, err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_ZeroPrice(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)

	req := &domain.CreateProductRequest{
		Name:  "Test Product",
		Price: 0,
	}

	result, err := svc.CreateProduct(context.Background(), req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, constant.ErrInvalidPrice, err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_NegativePrice(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)

	req := &domain.CreateProductRequest{
		Name:  "Test Product",
		Price: -50.00,
	}

	result, err := svc.CreateProduct(context.Background(), req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, constant.ErrInvalidPrice, err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_SalePriceGreaterThanPrice(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)

	salePrice := 150.00
	req := &domain.CreateProductRequest{
		Name:      "Test Product",
		SalePrice: &salePrice,
		Price:     100.00,
	}

	result, err := svc.CreateProduct(context.Background(), req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, constant.ErrInvalidSalePrice, err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_SalePriceEqualToPrice(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)

	salePrice := 100.00
	req := &domain.CreateProductRequest{
		Name:      "Test Product",
		SalePrice: &salePrice,
		Price:     100.00,
	}

	result, err := svc.CreateProduct(context.Background(), req)

	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, constant.ErrInvalidSalePrice, err.Error())
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUpdateProduct_Success_PartialUpdate(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)

	existing := &domain.Product{
		ID:    1,
		Name:  "Old Name",
		Price: 100.00,
	}

	mockRepo.On("GetByID", mock.Anything, 1).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *domain.Product) bool {
		return p.Name == "New Name" && p.Price == 100.00
	})).Return(nil)

	req := &domain.UpdateProductRequest{
		Name: json.RawMessage(`"New Name"`),
	}
	err := svc.UpdateProduct(context.Background(), 1, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Success_SetNull(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)

	desc := "Old Desc"
	existing := &domain.Product{
		ID:          1,
		Name:        "Name",
		Description: &desc,
		Price:       100.00,
	}

	mockRepo.On("GetByID", mock.Anything, 1).Return(existing, nil)
	mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(p *domain.Product) bool {
		return p.Description == nil
	})).Return(nil)

	req := &domain.UpdateProductRequest{
		Description: json.RawMessage(`null`),
	}
	err := svc.UpdateProduct(context.Background(), 1, req)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	mockRepo := new(mocks.MockProductRepository)
	svc := NewProductService(mockRepo)

	mockRepo.On("GetByID", mock.Anything, 999).Return(nil, sql.ErrNoRows)

	req := &domain.UpdateProductRequest{
		Name: json.RawMessage(`"New Name"`),
	}
	err := svc.UpdateProduct(context.Background(), 999, req)

	assert.Error(t, err)
	assert.Equal(t, constant.ErrProductNotFound, err.Error())
	mockRepo.AssertNotCalled(t, "Update")
}



