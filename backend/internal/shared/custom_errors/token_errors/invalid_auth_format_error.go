package token_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	InvalidAuthFormatErrorCode abstract_error_code.ErrorCode = "INVALID_AUTH_FORMAT"
)

type InvalidAuthFormatError struct {
	message string
}

func NewInvalidAuthFormatError() *InvalidAuthFormatError {
	return &InvalidAuthFormatError{message: "invalid authorization header format"}
}

func (e *InvalidAuthFormatError) Error() string {
	return e.message
}

func (e *InvalidAuthFormatError) Code() abstract_error_code.ErrorCode {
	return InvalidAuthFormatErrorCode
}

func (e *InvalidAuthFormatError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateUnauthorizedResponse(InvalidAuthFormatErrorCode, e.Error())
}
