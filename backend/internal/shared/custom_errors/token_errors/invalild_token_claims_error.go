package token_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	InvalidTokenClaimsErrorCode abstract_error_code.ErrorCode = "INVALID_TOKEN_CLAIMS"
)

type InvalidTokenClaimsError struct {
	message string
}

func NewInvalidTokenClaimsError() *InvalidTokenClaimsError {
	return &InvalidTokenClaimsError{message: "invalid token claims"}
}

func (e *InvalidTokenClaimsError) Error() string {
	return e.message
}

func (e *InvalidTokenClaimsError) Code() abstract_error_code.ErrorCode {
	return InvalidTokenClaimsErrorCode
}

func (e *InvalidTokenClaimsError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateUnauthorizedResponse(InvalidAuthFormatErrorCode, e.Error())
}
