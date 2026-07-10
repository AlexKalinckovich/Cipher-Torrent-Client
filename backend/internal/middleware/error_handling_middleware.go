package middleware

import (
	"errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
	"github.com/gin-gonic/gin"
)

func ErrorHandlingMiddleware(registry *transport.ErrorRegistry) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		handleErrorsIfPresent(c, registry)
	}
}

func handleErrorsIfPresent(c *gin.Context, registry *transport.ErrorRegistry) {
	if len(c.Errors) == 0 {
		return
	}
	writeErrorResponse(c, registry)
}

func writeErrorResponse(c *gin.Context, registry *transport.ErrorRegistry) {
	lastError := c.Errors.Last()
	httpResp := resolveHTTPResponse(lastError, registry)
	c.JSON(httpResp.StatusCode(), httpResp.Body)
}

func resolveHTTPResponse(err error, registry *transport.ErrorRegistry) transport.HTTPResponse {
	if provider := tryExtractProvider(err); provider != nil {
		return provider.Handle()
	}
	return registry.Translate(err)
}

func tryExtractProvider(err error) transport.ErrorHandlerProvider {
	var provider transport.ErrorHandlerProvider
	if errors.As(err, &provider) {
		return provider
	}
	return nil
}
