package token_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	InvalidTokenErrorCode abstract_error_code.ErrorCode = "INVALID_TOKEN"
)

type InvalidTokenError struct {
	message string
	cause   error
}

func NewInvalidTokenError(cause error) *InvalidTokenError {
	return &InvalidTokenError{message: "invalid or expired token", cause: cause}
}

func (e *InvalidTokenError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *InvalidTokenError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateUnauthorizedResponse(InvalidTokenErrorCode, e.Error())
}
