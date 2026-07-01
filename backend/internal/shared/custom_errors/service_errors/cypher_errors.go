package service_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
)

const (
	DpkiErrorCode       abstract_error_code.ErrorCode = "DPKI_KEY_GENERATION_FAILED"
	EncryptionErrorCode abstract_error_code.ErrorCode = "ENCRYPTION_FAILED"
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

type EncryptionError struct {
	cause error
}

func NewEncryptionError(cause error) *EncryptionError {
	return &EncryptionError{cause: cause}
}

func (e *EncryptionError) Code() abstract_error_code.ErrorCode { return EncryptionErrorCode }
func (e *EncryptionError) Error() string                       { return "failed to encrypt private key" }
func (e *EncryptionError) Unwrap() error                       { return e.cause }
