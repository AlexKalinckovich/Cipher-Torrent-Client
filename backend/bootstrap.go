package main

import (
	"database/sql"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/auth_handler"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/torrent_handler"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/torrent_signature_handler"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers/user_handler"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/minio"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent_signature"
	torrent_signing "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent_signature"
	signing2 "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent_signature/meta_info_signer"
	signing "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent_signature/signature_builder"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	torrentDb "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/generated"
	userDb "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/config"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/torrent/meta_info"
	userMapper "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/event_bus"
	eventPorts "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/token_store"
	torrentRepo "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/torrent"
	userRepo "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/auth"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/infra"
	userService "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/auth/jwt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/dpki"
	crypto "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/security"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/security/password_hasher"
	storagePorts "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/storage/ports"
	userValidator "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/validator/user"
)

func bootstrapModules(rg *gin.RouterGroup, conn *sql.DB, redisClient *redis.Client) {
	userSvc := bootstrapUserModule(rg, conn)
	bootstrapAuthModule(rg, redisClient, userSvc)
	bootstrapTorrentModule(rg, conn, redisClient)
	bootstrapTorrentSignatureModule(rg, conn)
}

func bootstrapUserModule(rg *gin.RouterGroup, conn *sql.DB) *userService.UserService {
	queries := userDb.New(conn)
	repository := userRepo.NewUserRepository(conn, queries)
	validator := userValidator.NewUserValidator()
	mapper := userMapper.NewUserDTOMapper()
	masterKey := config.GetRequiredEnv("AES_MASTER_KEY")
	cr := initializeCryptoService(masterKey)
	service := userService.NewUserService(repository, mapper, validator, cr)
	handler := user_handler.NewUserHandler(service)
	handler.RegisterRoutes(rg)
	return service
}

func bootstrapTorrentSignatureModule(rg *gin.RouterGroup, conn *sql.DB) {
	userQueries := userDb.New(conn)
	userRepository := userRepo.NewUserRepository(conn, userQueries)

	signatureQueries := torrentDb.New(conn)
	signatureRepository := torrent_signature.NewTorrentSignatureRepository(conn, signatureQueries)

	masterKey := config.GetRequiredEnv("AES_MASTER_KEY")
	cryptoService := initializeCryptoService(masterKey)

	signatureBuilder := signing.NewSignatureBuilder()
	metaInfoSigner := signing2.NewMetaInfoSigner()
	storageRepo := initializeMinioStorage()

	service := torrent_signing.NewService(
		userRepository,
		signatureRepository,
		cryptoService,
		signatureBuilder,
		storageRepo,
		metaInfoSigner,
	)

	handler := torrent_signature_handler.NewTorrentSignatureHandler(service)
	handler.RegisterRoutes(rg)
}

func bootstrapAuthModule(rg *gin.RouterGroup, redisClient *redis.Client, userSvc *userService.UserService) {
	hasher := password_hasher.NewBcryptHasher()
	jwtSecret := []byte(config.GetRequiredEnv("JWT_SECRET"))
	tokenIssuer := jwt.NewJWTIssuer(jwtSecret, 15*time.Minute)
	refreshStore := token_store.NewRefreshTokenStore(redisClient)
	authSvc := auth.NewAuthService(userSvc, hasher, tokenIssuer, refreshStore)
	handler := auth_handler.NewAuthHandler(authSvc)
	handler.RegisterRoutes(rg)
}

func bootstrapTorrentModule(rg *gin.RouterGroup, conn *sql.DB, redisClient *redis.Client) {
	redisPublisher := event_bus.NewEventBus(redisClient)
	keyPair := generateIdentity()
	engine := initializeTorrentEngine(keyPair.PublicKey, redisPublisher)
	storageRepo := initializeMinioStorage()
	torrentQueries := torrentDb.New(conn)
	repository := torrentRepo.NewTorrentRepository(conn, torrentQueries)
	userQueries := userDb.New(conn)
	userRepository := userRepo.NewUserRepository(conn, userQueries)
	mapper := meta_info.NewMetainfoMapper()
	service := torrent.NewService(repository, engine, storageRepo, userRepository, mapper)
	handler := torrent_handler.NewTorrentHandler(service)
	handler.RegisterRoutes(rg)
}

func initializeCryptoService(masterKey string) crypto.CryptoServicePort {
	cr, err := crypto.NewCryptoService(masterKey)
	config.HandleInitError(err)
	return cr
}

func initializeMinioStorage() storagePorts.TorrentStoragePort {
	endpoint := config.GetRequiredEnv("MINIO_ENDPOINT")
	accessKey := config.GetRequiredEnv("MINIO_ROOT_USER")
	secretKey := config.GetRequiredEnv("MINIO_ROOT_PASSWORD")
	bucketName := config.GetRequiredEnv("MINIO_BUCKET_NAME")
	repo, err := minio.NewMinioStorageRepository(endpoint, accessKey, secretKey, bucketName, false)
	config.HandleInitError(err)
	return repo
}

func initializeTorrentEngine(pubKey []byte, publisher eventPorts.EventPublisher) *infra.AnacrolixEngine {
	dataDir := config.GetRequiredEnv("TORRENT_DATA_DIR")
	engine, err := infra.NewAnacrolixEngine(dataDir, pubKey, publisher)
	config.HandleInitError(err)
	return engine
}

func generateIdentity() *dpki.KeyPair {
	keyPair, err := dpki.GenerateIdentity()
	config.HandleInitError(err)
	return keyPair
}
