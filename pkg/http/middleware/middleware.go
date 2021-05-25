package middleware

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func respondWithError(c *gin.Context, code int, message interface{}) {
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

func TokenAuthMiddleware(apiId, apiCode string) gin.HandlerFunc {
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", apiId, apiCode)))

	return func(c *gin.Context) {
		remote := c.Request.Header.Get("Authorization")

		if remote != token {
			respondWithError(c, http.StatusForbidden, "")
			return
		}

		c.Next()
	}
}
