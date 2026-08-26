package http

import (
	"net/http"

	"product-api-assessment/internal/domain"
	"product-api-assessment/internal/handler/dto"
	"product-api-assessment/pkg/constant"
	"product-api-assessment/pkg/response"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	service domain.ProductService
}

func NewProductHandler(service domain.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

// CreateProduct godoc
// @Summary      Create a new product
// @Description  Create a new product with the given details
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        request  body      dto.CreateProductRequest  true  "Product data"
// @Success      201      {object}  response.Response{data=domain.Product}
// @Failure      400      {object}  response.Response
// @Router       /v1/product [post]
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(constant.ErrInvalidRequest))
		return
	}
	domainReq := &domain.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		SalePrice:   req.SalePrice,
		Price:       req.Price,
	}

	product, err := h.service.CreateProduct(c.Request.Context(), domainReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success(product))
}
