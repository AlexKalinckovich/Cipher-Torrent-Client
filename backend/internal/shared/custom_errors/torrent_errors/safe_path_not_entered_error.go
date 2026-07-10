package torrent_errors

import (
	"net/http"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/common"
)

const SafePathNotSpecifiedErrorCode abstract_error_code.ErrorCode = "SAFE_PATH_NOT_SPECIFIED"

type SafePathNotSpecifiedError struct{}

func NewSafePathNotSpecifiedError() *SafePathNotSpecifiedError {
	return &SafePathNotSpecifiedError{}
}

func (e *SafePathNotSpecifiedError) Error() string {
	return "SafePath not specified in header"
}

func (e *SafePathNotSpecifiedError) Code() abstract_error_code.ErrorCode {
	return SafePathNotSpecifiedErrorCode
}

func (e *SafePathNotSpecifiedError) Handle() transport.HTTPResponse {
	return buildSafePathResponse()
}

func buildSafePathResponse() transport.HTTPResponse {
	return transport.NewHTTPResponse(http.StatusBadRequest, common.ApiError{
		Status:    http.StatusBadRequest,
		ErrorCode: string(SafePathNotSpecifiedErrorCode),
		Message:   "Enter safepath into header",
	})
}
