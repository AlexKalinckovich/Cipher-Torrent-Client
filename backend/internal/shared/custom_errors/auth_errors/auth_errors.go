package auth_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
)

const (
	InvalidCredentialsErrorCode abstract_error_code.ErrorCode = "INVALID_CREDENTIALS"
	TokenGenerationErrorCode    abstract_error_code.ErrorCode = "TOKEN_GENERATION_FAILED"
	RefreshTokenSaveErrorCode   abstract_error_code.ErrorCode = "REFRESH_TOKEN_SAVE_FAILED"
)

type InvalidCredentialsError struct {
	message string
}

func NewInvalidCredentialsError() *InvalidCredentialsError {
	return &InvalidCredentialsError{message: "invalid email or password"}
}

func (e *InvalidCredentialsError) Error() string {
	return e.message
}

func (e *InvalidCredentialsError) Code() abstract_error_code.ErrorCode {
	return InvalidCredentialsErrorCode
}

type TokenGenerationError struct {
	message string
	cause   error
}

func NewTokenGenerationError(cause error) *TokenGenerationError {
	return &TokenGenerationError{message: "failed to generate token", cause: cause}
}

func (e *TokenGenerationError) Error() string {
	if e.cause != nil {
		return e.message + ": " + e.cause.Error()
	}
	return e.message
}

func (e *TokenGenerationError) Code() abstract_error_code.ErrorCode {
	return TokenGenerationErrorCode
}

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
