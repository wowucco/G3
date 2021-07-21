package http

import (
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/contact"
)

func RegisterHTTPEndpoints(router *gin.RouterGroup, platformAuth gin.HandlerFunc, contactUC contact.IContactUseCase) {
	h := NewHandler(contactUC)

	c := router.Group("/contact")

	c.Use(platformAuth)

	{
		c.POST("recall", h.recall)
		c.POST("buy-on-click", h.buyOnClick)
	}
}
