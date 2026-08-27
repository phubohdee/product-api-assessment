package http

import (
	"net/http"
	"strconv"

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
// @Failure      500      {object}  response.Response
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
		if err.Error() == constant.ErrInternalServer {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, response.Success(product))
}

// UpdateProduct godoc
// @Summary      Update a product
// @Description  Update only the fields that are sent in the request body
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true  "Product ID"
// @Param        request  body      dto.UpdateProductRequest  true  "Fields to update"
// @Success      200      {object}  response.Response
// @Failure      400      {object}  response.Response
// @Failure      404      {object}  response.Response
// @Failure      500      {object}  response.Response
// @Router       /v1/product/{id} [patch]
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(constant.ErrInvalidRequest))
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(constant.ErrInvalidRequest))
		return
	}

	domainReq := &domain.UpdateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		SalePrice:   req.SalePrice,
		Price:       req.Price,
	}

	if err := h.service.UpdateProduct(c.Request.Context(), id, domainReq); err != nil {
		if err.Error() == constant.ErrProductNotFound {
			c.JSON(http.StatusNotFound, response.Error(err.Error()))
			return
		}
		if err.Error() == constant.ErrInternalServer {
			c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
			return
		}
		c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}
