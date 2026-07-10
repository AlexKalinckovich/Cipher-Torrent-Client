package routers

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/auth_handler"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/torrent_handler"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/user_handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/api/v1")

	torrentHandler := torrent_handler.NewTorrentHandler()
	userHandler := user_handler.NewUserHandler()
	authHandler := auth_handler.NewAuthHandler()

	torrentHandler.RegisterRoutes(v1)
	userHandler.RegisterRoutes(v1)
	authHandler.RegisterRoutes(v1)

	return router
}
