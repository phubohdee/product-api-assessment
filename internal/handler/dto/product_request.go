package dto

import "encoding/json"

type CreateProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Description *string  `json:"description"`
	SalePrice   *float64 `json:"sale_price"`
	Price       float64  `json:"price" binding:"required"`
}

type UpdateProductRequest struct {
	Name        json.RawMessage `json:"name" swagtype:"string"`
	Description json.RawMessage `json:"description" swagtype:"string"`
	SalePrice   json.RawMessage `json:"sale_price" swagtype:"number"`
	Price       json.RawMessage `json:"price" swagtype:"number"`
}
