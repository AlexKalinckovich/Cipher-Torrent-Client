package auth_errors

import "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"

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
