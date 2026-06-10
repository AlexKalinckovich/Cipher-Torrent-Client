package main

import (
	"errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/middleware"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/validation"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
	openapi "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func SetupRouter(r *gin.Engine) *gin.Engine {
	registry := transport.NewErrorRegistry()

	registry.Register(validation.AggregateErrorCode, func(err error) transport.HTTPResponse {
		var aggErr *validation.AggregateError
		if errors.As(err, &aggErr) {
			return transport.NewHTTPResponse(http.StatusBadRequest, map[string]any{
				"code":    aggErr.Code(),
				"message": "Validation failed",
				"details": aggErr.Errors,
			})
		}
		return transport.DefaultFallbackHandler(err)
	})

	r.Use(middleware.ErrorHandlingMiddleware(registry))

	r.POST("/api/v1/torrent/sign", func(c *gin.Context) {
		aggErr := validation.NewAggregateError()

		signature := c.GetHeader("X-Signature")
		if signature == "" {
			aggErr.Add("X-Signature", nil, "signature is missing")
		}

		if aggErr.HasErrors() {
			err := c.Error(aggErr)

			currErrString := err.Error()
			aggErrString := aggErr.Error()

			if !strings.EqualFold(currErrString, aggErrString) {
				log.Println("Different Error registered ", currErrString, aggErrString)
			}
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "torrent signed successfully"})
	})

	return r
}

func main() {
	routes := openapi.ApiHandleFunctions{}
	log.Printf("Server started")
	router := openapi.NewRouter(routes)
	SetupRouter(router)
	log.Fatal(router.Run(":8080"))
}
