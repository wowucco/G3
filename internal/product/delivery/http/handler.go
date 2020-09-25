package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wowucco/G3/internal/product"
	"github.com/wowucco/G3/pkg/pagination"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase product.UseCase
}

func NewHandler(useCase product.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

func (h *Handler) get (c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	p, err := h.useCase.Get(c.Request.Context(), id)

	if err != nil {
		fmt.Print(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if p == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	c.JSON(http.StatusOK, &getResponse{
		Product: toProduct(p),
	})
}

func (h *Handler) all(c *gin.Context) {
	count, err := h.useCase.Count(c.Request.Context())

	if err != nil {
		fmt.Print(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	pages := pagination.NewFromRequest(c.Request, count)

	products, err := h.useCase.Query(c.Request.Context(), pages.Offset(), pages.Limit())

	if err != nil {
		fmt.Print(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	pages.Items = toProducts(products)

	c.JSON(http.StatusOK, gin.H{
		"products": pages,
	})
}

func (h *Handler) delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	p, err := h.useCase.Delete(c.Request.Context(), id)

	if p == nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusOK)
}