package middleware

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware(registry *transport.ErrorRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		isErrorOccurred := len(c.Errors) > 0

		if isErrorOccurred {
			handleError(c, registry)
		}
	}
}

func handleError(c *gin.Context, registry *transport.ErrorRegistry) {
	lastError := c.Errors.Last()

	httpResp := registry.Translate(lastError)

	c.JSON(httpResp.StatusCode(), httpResp.Body)
}
