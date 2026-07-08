package main

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/config"
	"net/http"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/middleware"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/common"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnvironment()
	dbConn := config.InitializeDatabase()
	defer dbConn.Close()

	engine := setupHTTPServer()
	errorRegistry := buildErrorRegistry()

	engine.Use(middleware.ErrorHandlingMiddleware(errorRegistry))
	engine.NoRoute(handleNoRoute)

	redisClient := config.InitializeRedisClient()

	v1 := engine.Group("/api/v1")

	bootstrapModules(v1, dbConn, redisClient)

	startServer(engine)
}

func loadEnvironment() {
	err := godotenv.Load("../containerization/.env")
	config.HandleInitError(err)
}

func setupHTTPServer() *gin.Engine {
	return gin.Default()
}

func handleNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, common.ApiError{
		Status:    http.StatusNotFound,
		ErrorCode: "ROUTE_NOT_FOUND",
		Message:   "the requested route does not exist",
	})
}

func startServer(engine *gin.Engine) {
	runErr := engine.Run(":8080")
	config.HandleInitError(runErr)
}
