package routers

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	torrentHandler := handlers.NewTorrentHandler()
	userHandler := handlers.NewUserHandler()
	authHandler := handlers.NewAuthHandler()

	torrentHandler.RegisterRoutes(v1)
	userHandler.RegisterRoutes(v1)
	authHandler.RegisterRoutes(v1)

	return router
}
