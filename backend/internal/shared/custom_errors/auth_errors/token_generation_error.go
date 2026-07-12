package auth_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const TokenGenerationErrorCode abstract_error_code.ErrorCode = "TOKEN_GENERATION_FAILED"

type TokenGenerationError struct {
	message string
	cause   error
}

func NewTokenGenerationError(cause error) *TokenGenerationError {
	return &TokenGenerationError{message: "failed to generate token", cause: cause}
}

func (e *TokenGenerationError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *TokenGenerationError) Code() abstract_error_code.ErrorCode {
	return TokenGenerationErrorCode
}

func (e *TokenGenerationError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateInternalResponse(TokenGenerationErrorCode, e.Error())
}
