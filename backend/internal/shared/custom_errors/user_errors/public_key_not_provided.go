package user_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const PublicKeyNotProvidedErrorCode abstract_error_code.ErrorCode = "PUBLIC_KEY_NOT_PROVIDED"

type PublicKeyNotProvidedError struct {
	message string
}

func NewPublicKeyNotProvidedError() *PublicKeyNotProvidedError {
	return &PublicKeyNotProvidedError{message: "provide public key as query param"}
}

func (e *PublicKeyNotProvidedError) Error() string {
	return e.message
}

func (e *PublicKeyNotProvidedError) Code() abstract_error_code.ErrorCode {
	return PublicKeyNotProvidedErrorCode
}

func (e *PublicKeyNotProvidedError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateConflictResponse(PublicKeyNotProvidedErrorCode, e.Error())
}
