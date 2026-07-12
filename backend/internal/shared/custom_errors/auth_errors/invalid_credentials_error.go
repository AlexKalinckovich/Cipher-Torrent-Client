package auth_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const InvalidCredentialsErrorCode abstract_error_code.ErrorCode = "INVALID_CREDENTIALS"

type InvalidCredentialsError struct {
	message string
}

func NewInvalidCredentialsError() *InvalidCredentialsError {
	return &InvalidCredentialsError{message: "invalid email or password"}
}

func (e *InvalidCredentialsError) Error() string {
	return e.message
}

func (e *InvalidCredentialsError) Code() abstract_error_code.ErrorCode {
	return InvalidCredentialsErrorCode
}

func (e *InvalidCredentialsError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateUnauthorizedResponse(InvalidRefreshTokenErrorCode, e.Error())
}
