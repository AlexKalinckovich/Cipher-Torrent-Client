package torrent_signature_repository_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const DuplicateSignatureErrorCode abstract_error_code.ErrorCode = "DUPLICATE_SIGNATURE"

type DuplicateSignatureError struct{}

func NewDuplicateSignatureError() *DuplicateSignatureError {
	return &DuplicateSignatureError{}
}

func (e *DuplicateSignatureError) Error() string {
	return "user has already signed this torrent"
}

func (e *DuplicateSignatureError) Code() abstract_error_code.ErrorCode {
	return DuplicateSignatureErrorCode
}

func (e *DuplicateSignatureError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateConflictResponse(DuplicateSignatureErrorCode, e.Error())
}
