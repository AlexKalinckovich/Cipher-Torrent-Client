package redis_errors

import (
	"fmt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const RedisFailedConnectionErrorCode = "RedisFailedConnectionError"

type RedisFailedConnectionError struct {
	message string
	cause   error
}

func NewRedisFailedConnectionError(message string) *RedisFailedConnectionError {
	return &RedisFailedConnectionError{message: message}
}
func NewRedisFailedConnectionErrorWithCause(message string, cause error) *RedisFailedConnectionError {
	return &RedisFailedConnectionError{
		message: message,
		cause:   cause,
	}
}

func (err *RedisFailedConnectionError) Error() string {
	cause := err.cause
	if cause != nil {
		return fmt.Sprintf("Redis FailedConnectionError.%s: %s", err.message, cause.Error())
	}
	return fmt.Sprintf("Redis FailedConnectionError.%s", err.message)
}

func (err *RedisFailedConnectionError) Code() string {
	return RedisFailedConnectionErrorCode
}

func (err *RedisFailedConnectionError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateInternalResponse(RedisFailedConnectionErrorCode, err.Error())
}
