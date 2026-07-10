package crypto_errors

import "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"

const (
	MasterKeyErrorCode abstract_error_code.ErrorCode = "INVALID_MASTER_KEY"
)

type MasterKeyError struct{}

func (e *MasterKeyError) Code() abstract_error_code.ErrorCode {
	return MasterKeyErrorCode
}

func (e *MasterKeyError) Error() string {
	return "invalid master key format or length"
}
