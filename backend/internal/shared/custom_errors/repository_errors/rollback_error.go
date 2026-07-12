package repository_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const RollbackErrorCode abstract_error_code.ErrorCode = "REPOSITORY_ROLLBACK_ERROR"

type RollbackError struct {
	OriginalErr error
	RollbackErr error
}

func NewRollbackError(originalErr error, rollbackErr error) *RollbackError {
	return &RollbackError{
		OriginalErr: originalErr,
		RollbackErr: rollbackErr,
	}
}

func (e *RollbackError) Code() abstract_error_code.ErrorCode {
	return RollbackErrorCode
}

func (e *RollbackError) Error() string {
	return "rollback failed: " + e.RollbackErr.Error() + "; original error: " + e.OriginalErr.Error()
}

func (e *RollbackError) Unwrap() error {
	return e.OriginalErr
}

func (e *RollbackError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateInternalResponse(RollbackErrorCode, e.Error())
}