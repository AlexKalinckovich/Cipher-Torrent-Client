package crypto_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	DpkiErrorCode abstract_error_code.ErrorCode = "DPKI_KEY_GENERATION_FAILED"
)

type DpkiError struct {
	cause error
}

func NewDpkiError(cause error) *DpkiError {
	return &DpkiError{cause: cause}
}

func (e *DpkiError) Code() abstract_error_code.ErrorCode { return DpkiErrorCode }
func (e *DpkiError) Error() string                       { return "failed to generate identity keypair" }
func (e *DpkiError) Unwrap() error                       { return e.cause }
func (e *DpkiError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateInternalResponse(DpkiErrorCode, e.Error())
}
