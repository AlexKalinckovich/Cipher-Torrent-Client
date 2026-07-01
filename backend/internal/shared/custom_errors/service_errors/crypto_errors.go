package service_errors

import "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"

const (
	MasterKeyErrorCode  abstract_error_code.ErrorCode = "INVALID_MASTER_KEY"
	DecryptionErrorCode abstract_error_code.ErrorCode = "DECRYPTION_FAILED"
)

type MasterKeyError struct{}

func (e *MasterKeyError) Code() abstract_error_code.ErrorCode {
	return MasterKeyErrorCode
}

func (e *MasterKeyError) Error() string {
	return "invalid master key format or length"
}

type DecryptionError struct {
	cause error
}

func NewDecryptionError(cause error) *DecryptionError {
	return &DecryptionError{cause: cause}
}

func (e *DecryptionError) Code() abstract_error_code.ErrorCode {
	return DecryptionErrorCode
}

func (e *DecryptionError) Error() string {
	return "failed to decrypt private key"
}

func (e *DecryptionError) Unwrap() error {
	return e.cause
}
