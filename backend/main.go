package main

import (
	"database/sql"
	"errors"
	"fmt"
	torrent2 "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/torrent"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/redis_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/infra"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/not_found"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers"
	userDtoMapper "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/middleware"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user"
	repositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user/repository_errors"
	userService "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	db "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
	sharedRepositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/repository_errors"
	serviceErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/service_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/validation"
	crypto "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/security"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
	userValidator "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/validator/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/common"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../containerization/.env")
	handleInitError(err)

	conn, connErr := sql.Open("mysql", buildDSN())
	handleInitError(connErr)
	defer conn.Close()

	configureConnectionPool(conn)

	engine := gin.Default()
	registry := buildErrorRegistry()
	engine.Use(middleware.ErrorHandlingMiddleware(registry))
	engine.NoRoute(handleNoRoute)

	v1 := engine.Group("/api/v1")
	bootstrapUserModule(v1, conn)
	bootstrapTorrentModule(v1)

	runErr := engine.Run(":8080")
	handleInitError(runErr)
}

func handleNoRoute(c *gin.Context) {
	c.JSON(http.StatusNotFound, common.ApiError{
		Status:    http.StatusNotFound,
		ErrorCode: "ROUTE_NOT_FOUND",
		Message:   "the requested route does not exist",
	})
}

func buildErrorRegistry() *transport.ErrorRegistry {
	registry := transport.NewErrorRegistry()
	registerValidationHandler(registry)
	registerRollbackHandler(registry)
	registerUserNotFoundHandler(registry)
	registerDuplicateEmailHandler(registry)
	registerDuplicateNicknameHandler(registry)
	registerDatabaseUnavailableHandler(registry)
	registerDpkiErrorHandler(registry)
	registerEncryptionErrorHandler(registry)
	registerRedisFailedConnectionError(registry)
	return registry
}

func registerRedisFailedConnectionError(registry *transport.ErrorRegistry) {
	registry.Register(redis_errors.RedisFailedConnectionErrorCode, func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusInternalServerError, common.ApiError{
			Status:    http.StatusInternalServerError,
			ErrorCode: redis_errors.RedisFailedConnectionErrorCode,
			Message:   "redis failed",
		})
	})
}

func registerValidationHandler(registry *transport.ErrorRegistry) {
	registry.Register(validation.AggregateErrorCode, func(err error) transport.HTTPResponse {
		var aggErr *validation.AggregateError
		if errors.As(err, &aggErr) {
			return transport.NewHTTPResponse(http.StatusBadRequest, common.ApiError{
				Status:    http.StatusBadRequest,
				ErrorCode: string(validation.AggregateErrorCode),
				Message:   "validation failed",
				Details:   buildValidationDetails(aggErr),
			})
		}
		return transport.DefaultFallbackHandler(err)
	})
}

func registerRollbackHandler(registry *transport.ErrorRegistry) {
	registry.Register(sharedRepositoryErrors.RollbackErrorCode, func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusInternalServerError, common.ApiError{
			Status:    http.StatusInternalServerError,
			ErrorCode: string(sharedRepositoryErrors.RollbackErrorCode),
			Message:   "a database transaction could not be completed",
		})
	})
}

func registerUserNotFoundHandler(registry *transport.ErrorRegistry) {
	code := not_found.CodeFor[repositoryErrors.UserEntity]()
	registry.Register(code, func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusNotFound, common.ApiError{
			Status:    http.StatusNotFound,
			ErrorCode: string(code),
			Message:   "the requested user does not exist",
		})
	})
}

func registerDuplicateEmailHandler(registry *transport.ErrorRegistry) {
	registry.Register(repositoryErrors.DuplicateEmailErrorCode, func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusConflict, common.ApiError{
			Status:    http.StatusConflict,
			ErrorCode: string(repositoryErrors.DuplicateEmailErrorCode),
			Message:   "a user with this email already exists",
		})
	})
}

func registerDuplicateNicknameHandler(registry *transport.ErrorRegistry) {
	registry.Register(repositoryErrors.DuplicateNicknameErrorCode, func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusConflict, common.ApiError{
			Status:    http.StatusConflict,
			ErrorCode: string(repositoryErrors.DuplicateNicknameErrorCode),
			Message:   "a user with this nickname already exists",
		})
	})
}

func registerDatabaseUnavailableHandler(registry *transport.ErrorRegistry) {
	registry.Register(repositoryErrors.DatabaseUnavailableCode, func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusServiceUnavailable, common.ApiError{
			Status:    http.StatusServiceUnavailable,
			ErrorCode: string(repositoryErrors.DatabaseUnavailableCode),
			Message:   "service temporarily unavailable",
		})
	})
}

func registerDpkiErrorHandler(registry *transport.ErrorRegistry) {
	registry.Register(serviceErrors.DpkiErrorCode, func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusInternalServerError, common.ApiError{
			Status:    http.StatusInternalServerError,
			ErrorCode: string(serviceErrors.DpkiErrorCode),
			Message:   "identity key generation failed",
		})
	})
}

func registerEncryptionErrorHandler(registry *transport.ErrorRegistry) {
	registry.Register(serviceErrors.EncryptionErrorCode, func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusInternalServerError, common.ApiError{
			Status:    http.StatusInternalServerError,
			ErrorCode: string(serviceErrors.EncryptionErrorCode),
			Message:   "key encryption failed",
		})
	})
}

func buildValidationDetails(aggErr *validation.AggregateError) map[string][]string {
	details := make(map[string][]string, len(aggErr.Errors))
	for _, fe := range aggErr.Errors {
		details[fe.Field] = append(details[fe.Field], fe.Message)
	}
	return details
}

func handleInitError(err error) {
	if err != nil {
		panic(err)
	}
}

func buildDSN() string {
	return fmt.Sprintf("root:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_ROOT_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_DATABASE"),
	)
}

func configureConnectionPool(conn *sql.DB) {
	conn.SetMaxOpenConns(10)
	conn.SetMaxIdleConns(5)
	conn.SetConnMaxLifetime(15 * time.Minute)
}

func bootstrapUserModule(rg *gin.RouterGroup, conn *sql.DB) {
	queries := db.New(conn)
	repository := user.NewUserRepository(conn, queries)
	validator := userValidator.NewUserValidator()
	mapper := userDtoMapper.NewUserDTOMapper()

	masterKey := os.Getenv("AES_MASTER_KEY")
	if masterKey == "" {
		panic("AES_MASTER_KEY environment variable is not set")
	}

	cr, err := crypto.NewCryptoService(masterKey)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize crypto service: %v", err))
	}

	service := userService.NewUserService(repository, mapper, validator, cr)
	handler := handlers.NewUserHandler(service)
	handler.RegisterRoutes(rg)
}

func bootstrapTorrentModule(rg *gin.RouterGroup) {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisPublisher, err := redis.NewEventBus(redisAddr, redisPassword)
	if err != nil {
		log.Fatalf("Failed to initialize Redis EventBus: %v", err)
	}

	engine, err := infra.NewAnacrolixEngine("./test_downloads", redisPublisher)
	if err != nil {
		log.Fatalf("Failed to initialize Torrent Engine: %v", err)
	}

	mapper := &torrent2.MetainfoMapper{}
	torrentService := torrent.NewService(engine, mapper)

	handler := handlers.NewTorrentHandler(torrentService)
	handler.RegisterRoutes(rg)
}
