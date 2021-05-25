package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/checkout"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, ordUC checkout.IOrderUseCase, platformAuth gin.HandlerFunc) {
	h := NewHandler(ordUC)

	c := router.Group("/checkout")
	c.Use(platformAuth)
	{
		c.POST("create", h.create)
		c.POST("init-payment", h.initPayment)
		c.POST("accept-holden-payment", h.acceptHolden)
	}

	cb := router.Group("/callback")
	{
		cb.POST(":provider", h.callback)
	}
}
