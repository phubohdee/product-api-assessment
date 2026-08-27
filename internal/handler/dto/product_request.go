package dto

type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required" example:"Product Name"`
	Description *string  `json:"description" example:"Product Description"`
	SalePrice   *float64 `json:"sale_price" example:"80.00"`
	Price       float64  `json:"price" binding:"required" example:"100.00"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name" example:"Updated Name"`
	Description *string  `json:"description" example:"Updated Description"`
	SalePrice   *float64 `json:"sale_price" example:"80.00"`
	Price       *float64 `json:"price" example:"100.00"`
}
