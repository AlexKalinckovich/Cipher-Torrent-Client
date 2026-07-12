package token_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	MissingAuthHeaderErrorCode abstract_error_code.ErrorCode = "MISSING_AUTH_HEADER"
)

type MissingAuthHeaderError struct {
	message string
}

func NewMissingAuthHeaderError() *MissingAuthHeaderError {
	return &MissingAuthHeaderError{message: "missing authorization header"}
}

func (e *MissingAuthHeaderError) Error() string {
	return e.message
}

func (e *MissingAuthHeaderError) Code() abstract_error_code.ErrorCode {
	return MissingAuthHeaderErrorCode
}

func (e *MissingAuthHeaderError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateUnauthorizedResponse(MissingAuthHeaderErrorCode, e.Error())
}
