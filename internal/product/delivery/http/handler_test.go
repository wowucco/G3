package http

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/wowucco/G3/internal/product"
	"github.com/wowucco/G3/internal/product/usecase"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDelete(t *testing.T) {
	r := gin.Default()
	group := r.Group("/api")

	uc := new(usecase.ProductUseCaseMock)

	RegisterHTTPEndpoints(group, uc)

	form := product.DeleteProductFrom{
		ID: 1234,
	}

	body, err := json.Marshal(form)
	assert.NoError(t, err)

	uc.On("Delete", form).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/products", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}