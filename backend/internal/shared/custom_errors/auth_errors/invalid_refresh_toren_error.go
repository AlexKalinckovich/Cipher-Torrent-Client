package auth_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const InvalidRefreshTokenErrorCode abstract_error_code.ErrorCode = "INVALID_REFRESH_TOKEN"

type InvalidRefreshTokenError struct {
	message string
}

func NewInvalidRefreshTokenError() *InvalidRefreshTokenError {
	return &InvalidRefreshTokenError{message: "invalid or expired refresh token"}
}

func (e *InvalidRefreshTokenError) Error() string {
	return e.message
}

func (e *InvalidRefreshTokenError) Code() abstract_error_code.ErrorCode {
	return InvalidRefreshTokenErrorCode
}

func (e *InvalidRefreshTokenError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateUnauthorizedResponse(InvalidRefreshTokenErrorCode, e.Error())
}
