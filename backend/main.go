package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/config"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/middleware"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/auth/jwt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/common"
)

func main() {
	loadEnvironment()
	dbConn := config.InitializeDatabase()
	defer dbConn.Close()
	engine := setupHTTPServer()
	registerGlobalMiddleware(engine)
	redisClient := config.InitializeRedisClient()
	registerV1Routes(engine, dbConn, redisClient)
	startServer(engine)
}

func loadEnvironment() {
	err := godotenv.Load("../containerization/.env")
	config.HandleInitError(err)
}

func setupHTTPServer() *gin.Engine {
	return gin.Default()
}

func registerGlobalMiddleware(engine *gin.Engine) {
	errorRegistry := buildErrorRegistry()
	engine.Use(middleware.ErrorHandlingMiddleware(errorRegistry))
	engine.NoRoute(handleNoRoute)
}

func handleNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, common.ApiError{
		Status:    http.StatusNotFound,
		ErrorCode: "ROUTE_NOT_FOUND",
		Message:   "the requested route does not exist",
	})
}

func registerV1Routes(engine *gin.Engine, dbConn *sql.DB, redisClient *redis.Client) {
	v1 := engine.Group("/api/v1")
	public := createPublicGroup(v1)
	private := createPrivateGroup(v1)
	userSvc := bootstrapUserModule(private, dbConn)
	bootstrapAuthModule(public, redisClient, userSvc)
	bootstrapTorrentModule(private, dbConn, redisClient)
	bootstrapTorrentSignatureModule(private, dbConn)
}

func createPublicGroup(v1 *gin.RouterGroup) *gin.RouterGroup {
	return v1.Group("")
}

func createPrivateGroup(v1 *gin.RouterGroup) *gin.RouterGroup {
	private := v1.Group("")
	private.Use(initializeAuthMiddleware().Handle())
	return private
}

func initializeAuthMiddleware() *middleware.AuthMiddleware {
	secret := []byte(config.GetRequiredEnv("JWT_SECRET"))
	parser := jwt.NewJWTParser(secret)
	return middleware.NewAuthMiddleware(parser)
}

func startServer(engine *gin.Engine) {
	runErr := engine.Run(":8080")
	config.HandleInitError(runErr)
}
