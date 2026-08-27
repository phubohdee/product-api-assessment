package dto

type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price" binding:"required"`
}
