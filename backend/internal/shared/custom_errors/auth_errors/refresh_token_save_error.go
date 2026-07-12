package auth_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	RefreshTokenSaveErrorCode abstract_error_code.ErrorCode = "REFRESH_TOKEN_SAVE_FAILED"
)

type RefreshTokenSaveError struct {
	message string
	cause   error
}

func NewRefreshTokenSaveError(cause error) *RefreshTokenSaveError {
	return &RefreshTokenSaveError{message: "failed to save refresh token", cause: cause}
}

func (e *RefreshTokenSaveError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *RefreshTokenSaveError) Code() abstract_error_code.ErrorCode {
	return RefreshTokenSaveErrorCode
}

func (e *RefreshTokenSaveError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateInternalResponse(RefreshTokenSaveErrorCode, e.Error())
}
