package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"product-api-assessment/internal/domain"
	"product-api-assessment/mocks"
	"product-api-assessment/pkg/constant"
	"product-api-assessment/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupRouter(mockService *mocks.MockProductService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	handler := NewProductHandler(mockService)
	r := gin.Default()
	r.POST("/v1/product", handler.CreateProduct)
	return r
}

func TestCreateProduct_Success(t *testing.T) {
	mockService := new(mocks.MockProductService)

	desc := "A test product"
	expectedProduct := &domain.Product{
		ID:          1,
		Name:        "Test Product",
		Description: &desc,
		Price:       100.50,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("CreateProduct", mock.Anything, mock.AnythingOfType("*domain.CreateProductRequest")).Return(expectedProduct, nil)

	r := setupRouter(mockService)

	body := `{"name":"Test Product","description":"A test product","price":100.50}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/product", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Successful)
	assert.Empty(t, resp.ErrorCode)
	assert.NotNil(t, resp.Data)
	mockService.AssertExpectations(t)
}

func TestCreateProduct_InvalidJSON(t *testing.T) {
	mockService := new(mocks.MockProductService)
	r := setupRouter(mockService)

	body := `{invalid json}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/product", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Successful)
	assert.Equal(t, constant.ErrInvalidRequest, resp.ErrorCode)
	mockService.AssertNotCalled(t, "CreateProduct")
}

func TestCreateProduct_MissingRequiredFields(t *testing.T) {
	mockService := new(mocks.MockProductService)
	r := setupRouter(mockService)

	body := `{}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/product", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Successful)
	assert.Equal(t, constant.ErrInvalidRequest, resp.ErrorCode)
	mockService.AssertNotCalled(t, "CreateProduct")
}

func TestCreateProduct_ServiceError(t *testing.T) {
	mockService := new(mocks.MockProductService)

	mockService.On("CreateProduct", mock.Anything, mock.AnythingOfType("*domain.CreateProductRequest")).
		Return(nil, errors.New(constant.ErrInvalidPrice))

	r := setupRouter(mockService)

	body := `{"name":"Test Product","price":100.50}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/product", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.False(t, resp.Successful)
	assert.Equal(t, constant.ErrInvalidPrice, resp.ErrorCode)
	mockService.AssertExpectations(t)
}

func TestCreateProduct_WithNullableFields(t *testing.T) {
	mockService := new(mocks.MockProductService)

	salePrice := 80.00
	expectedProduct := &domain.Product{
		ID:        1,
		Name:      "Sale Product",
		SalePrice: &salePrice,
		Price:     100.00,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("CreateProduct", mock.Anything, mock.AnythingOfType("*domain.CreateProductRequest")).Return(expectedProduct, nil)

	r := setupRouter(mockService)

	body := `{"name":"Sale Product","sale_price":80.00,"price":100.00}`
	req, _ := http.NewRequest(http.MethodPost, "/v1/product", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp response.Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.True(t, resp.Successful)
	assert.NotNil(t, resp.Data)
	mockService.AssertExpectations(t)
}
