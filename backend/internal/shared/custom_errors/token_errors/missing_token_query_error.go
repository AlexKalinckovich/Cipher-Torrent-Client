package token_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	MissingTokenQueryErrorCode abstract_error_code.ErrorCode = "MISSING_TOKEN_QUERY"
)

type MissingTokenQueryError struct {
	message string
}

func NewMissingTokenQueryError() *MissingTokenQueryError {
	return &MissingTokenQueryError{message: "missing token query parameter"}
}

func (e *MissingTokenQueryError) Error() string {
	return e.message
}

func (e *MissingTokenQueryError) Code() abstract_error_code.ErrorCode {
	return MissingTokenQueryErrorCode
}

func (e *MissingTokenQueryError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateUnauthorizedResponse(MissingTokenQueryErrorCode, e.Error())
}
