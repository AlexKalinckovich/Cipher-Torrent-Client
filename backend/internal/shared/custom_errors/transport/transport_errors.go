package transport

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
)

const (
	UnknownErrorCode  abstract_error_code.ErrorCode = "UNKNOWN_ERROR"
	InternalErrorCode abstract_error_code.ErrorCode = "INTERNAL_ERROR"
)

type UnknownError struct {
	message string
	cause   error
}

func NewUnknownError(message string) *UnknownError {
	return &UnknownError{message: message}
}

func NewUnknownErrorWithCause(message string, cause error) *UnknownError {
	return &UnknownError{message: message, cause: cause}
}

func (e *UnknownError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *UnknownError) Code() abstract_error_code.ErrorCode {
	return UnknownErrorCode
}

type InternalError struct {
	message string
}

func NewInternalError(message string) *InternalError {
	return &InternalError{message: message}
}

func (e *InternalError) Error() string {
	return e.message
}

func (e *InternalError) Code() abstract_error_code.ErrorCode {
	return InternalErrorCode
}
