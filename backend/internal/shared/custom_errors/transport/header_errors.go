package transport

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
)

const (
	MissingAuthHeaderErrorCode  abstract_error_code.ErrorCode = "MISSING_AUTH_HEADER"
	InvalidAuthFormatErrorCode  abstract_error_code.ErrorCode = "INVALID_AUTH_FORMAT"
	InvalidTokenErrorCode       abstract_error_code.ErrorCode = "INVALID_TOKEN"
	InvalidTokenClaimsErrorCode abstract_error_code.ErrorCode = "INVALID_TOKEN_CLAIMS"
	MissingTokenQueryErrorCode  abstract_error_code.ErrorCode = "MISSING_TOKEN_QUERY"
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

type InvalidAuthFormatError struct {
	message string
}

func NewInvalidAuthFormatError() *InvalidAuthFormatError {
	return &InvalidAuthFormatError{message: "invalid authorization header format"}
}

func (e *InvalidAuthFormatError) Error() string {
	return e.message
}

func (e *InvalidAuthFormatError) Code() abstract_error_code.ErrorCode {
	return InvalidAuthFormatErrorCode
}

type InvalidTokenError struct {
	message string
	cause   error
}

func NewInvalidTokenError(cause error) *InvalidTokenError {
	return &InvalidTokenError{message: "invalid or expired token", cause: cause}
}

func (e *InvalidTokenError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *InvalidTokenError) Code() abstract_error_code.ErrorCode {
	return InvalidTokenErrorCode
}

type InvalidTokenClaimsError struct {
	message string
}

func NewInvalidTokenClaimsError() *InvalidTokenClaimsError {
	return &InvalidTokenClaimsError{message: "invalid token claims"}
}

func (e *InvalidTokenClaimsError) Error() string {
	return e.message
}

func (e *InvalidTokenClaimsError) Code() abstract_error_code.ErrorCode {
	return InvalidTokenClaimsErrorCode
}

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
