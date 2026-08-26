package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "product-api-assessment/docs"
)

func NewRouter(productHandler *ProductHandler) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/v1")
	{
		v1.GET("/api-docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		v1.POST("/product", productHandler.CreateProduct)
	}

	return r
}
