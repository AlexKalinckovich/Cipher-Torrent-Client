package repository_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	DuplicateEmailErrorCode abstract_error_code.ErrorCode = "EMAIL_ALREADY_EXISTS"
)

type DuplicateEmailError struct{}

func (e *DuplicateEmailError) Code() abstract_error_code.ErrorCode { return DuplicateEmailErrorCode }
func (e *DuplicateEmailError) Error() string                       { return "email already exists" }
func (e *DuplicateEmailError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateConflictResponse(DuplicateEmailErrorCode, e.Error())
}
