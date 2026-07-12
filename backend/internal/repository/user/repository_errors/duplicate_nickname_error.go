package repository_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	DuplicateNicknameErrorCode abstract_error_code.ErrorCode = "NICKNAME_ALREADY_EXISTS"
)

type DuplicateNicknameError struct{}

func (e *DuplicateNicknameError) Code() abstract_error_code.ErrorCode {
	return DuplicateNicknameErrorCode
}
func (e *DuplicateNicknameError) Error() string { return "nickname already exists" }
func (e *DuplicateNicknameError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateConflictResponse(DuplicateNicknameErrorCode, e.Error())
}
