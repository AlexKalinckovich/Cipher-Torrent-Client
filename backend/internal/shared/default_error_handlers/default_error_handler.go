package default_error_handlers

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/common"
	"net/http"
)

func CreateBadRequestResponse(code abstract_error_code.ErrorCode, message string) transport.HTTPResponse {
	return transport.NewHTTPResponse(http.StatusUnauthorized, common.ApiError{
		Status:    http.StatusBadRequest,
		ErrorCode: string(code),
		Message:   message,
	})
}

func CreateUnauthorizedResponse(code abstract_error_code.ErrorCode, message string) transport.HTTPResponse {
	return transport.NewHTTPResponse(http.StatusUnauthorized, common.ApiError{
		Status:    http.StatusUnauthorized,
		ErrorCode: string(code),
		Message:   message,
	})
}

func CreateInternalResponse(code abstract_error_code.ErrorCode, message string) transport.HTTPResponse {
	return transport.NewHTTPResponse(http.StatusInternalServerError, common.ApiError{
		Status:    http.StatusInternalServerError,
		ErrorCode: string(code),
		Message:   message,
	})
}

func CreateNotFoundResponse(code abstract_error_code.ErrorCode) transport.HTTPResponse {
	return transport.NewHTTPResponse(http.StatusNotFound, common.ApiError{
		Status:    http.StatusNotFound,
		ErrorCode: string(code),
		Message:   "the requested user does not exist",
	})
}

func CreateConflictResponse(code abstract_error_code.ErrorCode, message string) transport.HTTPResponse {
	return transport.NewHTTPResponse(http.StatusConflict, common.ApiError{
		Status:    http.StatusConflict,
		ErrorCode: string(code),
		Message:   message,
	})
}

func CreateServiceUnavailableResponse(code abstract_error_code.ErrorCode, message string) transport.HTTPResponse {
	return transport.NewHTTPResponse(http.StatusServiceUnavailable, common.ApiError{
		Status:    http.StatusServiceUnavailable,
		ErrorCode: string(code),
		Message:   message,
	})
}
