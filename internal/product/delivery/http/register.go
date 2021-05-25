package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/product"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, uc product.UseCase, platformAuth gin.HandlerFunc) {
	h := NewHandler(uc)

	products := router.Group("/products")
	products.Use(platformAuth)
	{
		products.POST(":id/", h.get)
		products.POST("/", h.all)
		//products.DELETE(":id", h.delete)
	}
}