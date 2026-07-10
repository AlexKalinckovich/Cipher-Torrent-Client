package main

import (
	"database/sql"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/auth_handler"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/torrent_handler"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/user_handler"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/websocket_handler"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/torrent/meta_info"
	userMapper "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/user"
	torrent2 "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent"
	t_validator "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/validator/torrent"
	"time"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/config"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/infrastructure/websocket"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/event_broker"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/event_bus"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/token_store"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/auth"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent"
	torrentDb "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/generated"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/infra"
	userService "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	userDb "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/auth/jwt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/dpki"
	crypto "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/security"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/security/password_hasher"
	userValidator "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/validator/user"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func bootstrapUserModule(rg *gin.RouterGroup, conn *sql.DB) *userService.UserService {
	queries := userDb.New(conn)
	repository := user.NewUserRepository(conn, queries)
	validator := userValidator.NewUserValidator()
	mapper := userMapper.NewUserDTOMapper()
	masterKey := config.GetRequiredEnv("AES_MASTER_KEY")
	cr := initializeCryptoService(masterKey)

	service := userService.NewUserService(repository, mapper, validator, cr)

	handler := user_handler.NewUserHandler(service)
	handler.RegisterRoutes(rg)

	return service
}

func bootstrapAuthModule(rg *gin.RouterGroup, redisClient *redis.Client, userSvc *userService.UserService) {
	hasher := password_hasher.NewBcryptHasher()
	jwtSecret := []byte(config.GetRequiredEnv("JWT_SECRET"))
	tokenIssuer := jwt.NewJWTIssuer(jwtSecret, 15*time.Minute)
	refreshStore := token_store.NewRefreshTokenStore(redisClient)

	authSvc := auth.NewAuthService(userSvc, hasher, tokenIssuer, refreshStore)

	authHandler := auth_handler.NewAuthHandler(authSvc)
	authHandler.RegisterRoutes(rg)
}

func initializeCryptoService(masterKey string) crypto.CryptoServicePort {
	cr, err := crypto.NewCryptoService(masterKey)
	if err != nil {
		panic("failed to initialize crypto service: " + err.Error())
	}
	return cr
}

func bootstrapTorrentModule(rg *gin.RouterGroup, conn *sql.DB, redisClient *redis.Client) {
	redisPublisher := event_bus.NewEventBus(redisClient)
	keyPair := generateIdentity()
	engine := initializeTorrentEngine(keyPair.PublicKey, redisPublisher)
	queries := torrentDb.New(conn)
	repository := torrent2.NewTorrentRepository(conn, queries)
	mapper := meta_info.NewMetainfoMapper()
	validator := t_validator.NewTorrentValidator()
	service := torrent.NewService(repository, engine, mapper, validator)
	handler := torrent_handler.NewTorrentHandler(service)
	handler.RegisterRoutes(rg)
	broker := event_broker.NewEventBroker(redisClient)
	hub := websocket.NewHub(broker)
	wsHandler := websocket_handler.NewWebSocketHandler(hub)
	rg.GET("/ws/:infoHash", wsHandler.Handle)
}

func generateIdentity() *dpki.KeyPair {
	keyPair, err := dpki.GenerateIdentity()
	if err != nil {
		panic("failed to generate identity: " + err.Error())
	}
	return keyPair
}

func initializeTorrentEngine(pubKey []byte, publisher ports.EventPublisher) *infra.AnacrolixEngine {
	engine, err := infra.NewAnacrolixEngine("./test_downloads", pubKey, publisher)
	if err != nil {
		panic("failed to initialize Torrent Engine: " + err.Error())
	}
	return engine
}
