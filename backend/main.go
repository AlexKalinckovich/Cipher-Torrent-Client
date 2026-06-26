package main

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/handlers"
	userDtoMapper "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/mapper/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/middleware"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user"
	repositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user/repository_errors"
	userService "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	db "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user/generated"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/not_found"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/repository_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/validation"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
	userValidator "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/validator/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/common"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
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

	v1 := engine.Group("/api/v1")
	bootstrapUserModule(v1, conn)

	runErr := engine.Run(":8080")
	handleInitError(runErr)
}

func buildErrorRegistry() *transport.ErrorRegistry {
	registry := transport.NewErrorRegistry()
	registerValidationHandler(registry)
	registerRollbackHandler(registry)
	registerUserNotFoundHandler(registry)
	return registry
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

func buildValidationDetails(aggErr *validation.AggregateError) map[string][]string {
	details := make(map[string][]string, len(aggErr.Errors))
	for _, fe := range aggErr.Errors {
		details[fe.Field] = append(details[fe.Field], fe.Message)
	}
	return details
}

func registerRollbackHandler(registry *transport.ErrorRegistry) {
	registry.Register(repository_errors.RollbackErrorCode, func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusInternalServerError, common.ApiError{
			Status:    http.StatusInternalServerError,
			ErrorCode: string(repository_errors.RollbackErrorCode),
			Message:   "a database transaction could not be completed",
		})
	})
}

func registerUserNotFoundHandler(registry *transport.ErrorRegistry) {
	registry.Register(not_found.CodeFor[repositoryErrors.UserEntity](), func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusNotFound, common.ApiError{
			Status:    http.StatusNotFound,
			ErrorCode: string(not_found.CodeFor[repositoryErrors.UserEntity]()),
			Message:   "the requested user does not exist",
		})
	})
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
	service := userService.NewUserService(repository, mapper, validator)
	handler := handlers.NewUserHandler(service)
	handler.RegisterRoutes(rg)
}
