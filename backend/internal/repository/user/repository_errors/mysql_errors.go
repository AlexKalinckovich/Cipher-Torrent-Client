package repository_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
)

const (
	DuplicateEmailErrorCode    abstract_error_code.ErrorCode = "EMAIL_ALREADY_EXISTS"
	DuplicateNicknameErrorCode abstract_error_code.ErrorCode = "NICKNAME_ALREADY_EXISTS"
	DatabaseUnavailableCode    abstract_error_code.ErrorCode = "DATABASE_UNAVAILABLE"
)

type DuplicateEmailError struct{}

func (e *DuplicateEmailError) Code() abstract_error_code.ErrorCode { return DuplicateEmailErrorCode }
func (e *DuplicateEmailError) Error() string                       { return "email already exists" }

type DuplicateNicknameError struct{}

func (e *DuplicateNicknameError) Code() abstract_error_code.ErrorCode {
	return DuplicateNicknameErrorCode
}
func (e *DuplicateNicknameError) Error() string { return "nickname already exists" }

type DatabaseUnavailableError struct{}

func (e *DatabaseUnavailableError) Code() abstract_error_code.ErrorCode {
	return DatabaseUnavailableCode
}
func (e *DatabaseUnavailableError) Error() string { return "database unavailable" }
