package main

import (
	"errors"
	repositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/repository/user/repository_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	tokenErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/auth_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/not_found"
	"net/http"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/redis_errors"
	sharedRepositoryErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/repository_errors"
	serviceErrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/service_errors"
	headererrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/transport"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/validation"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/common"
)

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
	registerAuthErrors(registry)
	registerInvalidRefreshToken(registry)
	return registry
}

func registerValidationHandler(registry *transport.ErrorRegistry) {
	registry.Register(validation.AggregateErrorCode, handleValidationError)
}

func handleValidationError(err error) transport.HTTPResponse {
	var aggErr *validation.AggregateError
	if errors.As(err, &aggErr) {
		return buildValidationErrorResponse(aggErr)
	}
	return transport.DefaultFallbackHandler(err)
}

func buildValidationErrorResponse(aggErr *validation.AggregateError) transport.HTTPResponse {
	return transport.NewHTTPResponse(http.StatusBadRequest, common.ApiError{
		Status:    http.StatusBadRequest,
		ErrorCode: string(validation.AggregateErrorCode),
		Message:   "validation failed",
		Details:   buildValidationDetails(aggErr),
	})
}

func buildValidationDetails(aggErr *validation.AggregateError) map[string][]string {
	details := make(map[string][]string, len(aggErr.Errors))
	for _, fe := range aggErr.Errors {
		details[fe.Field] = append(details[fe.Field], fe.Message)
	}
	return details
}

func registerInvalidRefreshToken(registry *transport.ErrorRegistry) {
	registry.Register(tokenErrors.InvalidRefreshTokenErrorCode, createUnauthorizedHandler(tokenErrors.InvalidRefreshTokenErrorCode))
}

func registerRollbackHandler(registry *transport.ErrorRegistry) {
	registry.Register(sharedRepositoryErrors.RollbackErrorCode, createInternalHandler(sharedRepositoryErrors.RollbackErrorCode, "a database transaction could not be completed"))
}

func registerUserNotFoundHandler(registry *transport.ErrorRegistry) {
	code := not_found.CodeFor[repositoryErrors.UserEntity]()
	registry.Register(code, createNotFoundHandler(code))
}

func registerDuplicateEmailHandler(registry *transport.ErrorRegistry) {
	registry.Register(repositoryErrors.DuplicateEmailErrorCode, createConflictHandler(repositoryErrors.DuplicateEmailErrorCode, "a user with this email already exists"))
}

func registerDuplicateNicknameHandler(registry *transport.ErrorRegistry) {
	registry.Register(repositoryErrors.DuplicateNicknameErrorCode, createConflictHandler(repositoryErrors.DuplicateNicknameErrorCode, "a user with this nickname already exists"))
}

func registerDatabaseUnavailableHandler(registry *transport.ErrorRegistry) {
	registry.Register(repositoryErrors.DatabaseUnavailableCode, createServiceUnavailableHandler(repositoryErrors.DatabaseUnavailableCode, "service temporarily unavailable"))
}

func registerDpkiErrorHandler(registry *transport.ErrorRegistry) {
	registry.Register(serviceErrors.DpkiErrorCode, createInternalHandler(serviceErrors.DpkiErrorCode, "identity key generation failed"))
}

func registerEncryptionErrorHandler(registry *transport.ErrorRegistry) {
	registry.Register(serviceErrors.EncryptionErrorCode, createInternalHandler(serviceErrors.EncryptionErrorCode, "key encryption failed"))
}

func registerRedisFailedConnectionError(registry *transport.ErrorRegistry) {
	registry.Register(redis_errors.RedisFailedConnectionErrorCode, createInternalHandler(redis_errors.RedisFailedConnectionErrorCode, "redis failed"))
}

func registerAuthErrors(registry *transport.ErrorRegistry) {
	registry.Register(tokenErrors.InvalidCredentialsErrorCode, createUnauthorizedHandler(tokenErrors.InvalidCredentialsErrorCode))
	registry.Register(tokenErrors.TokenGenerationErrorCode, createInternalHandler(tokenErrors.TokenGenerationErrorCode, "failed to generate authentication token"))
	registry.Register(tokenErrors.RefreshTokenSaveErrorCode, createInternalHandler(tokenErrors.RefreshTokenSaveErrorCode, "failed to save session data"))
	registry.Register(headererrors.MissingAuthHeaderErrorCode, createUnauthorizedHandler(headererrors.MissingAuthHeaderErrorCode))
	registry.Register(headererrors.InvalidAuthFormatErrorCode, createUnauthorizedHandler(headererrors.InvalidAuthFormatErrorCode))
	registry.Register(headererrors.InvalidTokenErrorCode, createUnauthorizedHandler(headererrors.InvalidTokenErrorCode))
	registry.Register(headererrors.InvalidTokenClaimsErrorCode, createUnauthorizedHandler(headererrors.InvalidTokenClaimsErrorCode))
	registry.Register(headererrors.MissingTokenQueryErrorCode, createUnauthorizedHandler(headererrors.MissingTokenQueryErrorCode))
}

func createUnauthorizedHandler(code abstract_error_code.ErrorCode) transport.ErrorHandler {
	return func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusUnauthorized, common.ApiError{
			Status:    http.StatusUnauthorized,
			ErrorCode: string(code),
			Message:   err.Error(),
		})
	}
}

func createInternalHandler(code abstract_error_code.ErrorCode, message string) transport.ErrorHandler {
	return func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusInternalServerError, common.ApiError{
			Status:    http.StatusInternalServerError,
			ErrorCode: string(code),
			Message:   message,
		})
	}
}

func createNotFoundHandler(code abstract_error_code.ErrorCode) transport.ErrorHandler {
	return func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusNotFound, common.ApiError{
			Status:    http.StatusNotFound,
			ErrorCode: string(code),
			Message:   "the requested user does not exist",
		})
	}
}

func createConflictHandler(code abstract_error_code.ErrorCode, message string) transport.ErrorHandler {
	return func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusConflict, common.ApiError{
			Status:    http.StatusConflict,
			ErrorCode: string(code),
			Message:   message,
		})
	}
}

func createServiceUnavailableHandler(code abstract_error_code.ErrorCode, message string) transport.ErrorHandler {
	return func(err error) transport.HTTPResponse {
		return transport.NewHTTPResponse(http.StatusServiceUnavailable, common.ApiError{
			Status:    http.StatusServiceUnavailable,
			ErrorCode: string(code),
			Message:   message,
		})
	}
}
