package crypto_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	EncryptionErrorCode abstract_error_code.ErrorCode = "ENCRYPTION_FAILED"
)

type EncryptionError struct {
	cause error
}

func NewEncryptionError(cause error) *EncryptionError {
	return &EncryptionError{cause: cause}
}

func (e *EncryptionError) Code() abstract_error_code.ErrorCode { return EncryptionErrorCode }
func (e *EncryptionError) Error() string                       { return "failed to encrypt private key" }
func (e *EncryptionError) Unwrap() error                       { return e.cause }
func (e *EncryptionError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateInternalResponse(EncryptionErrorCode, e.Error())
}
