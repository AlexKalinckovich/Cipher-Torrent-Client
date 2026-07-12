package torrent_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/common"
	"net/http"
)

const TorrentNotFoundInClientErrorCode abstract_error_code.ErrorCode = "TORRENT_NOT_FOUND_IN_CLIENT"

type TorrentNotFoundInClientError struct{}

func NewTorrentNotFoundInClientError() *TorrentNotFoundInClientError {
	return &TorrentNotFoundInClientError{}
}

func (e *TorrentNotFoundInClientError) Error() string {
	return "torrent not found in client"
}

func (e *TorrentNotFoundInClientError) Code() abstract_error_code.ErrorCode {
	return TorrentNotFoundInClientErrorCode
}

func (e *TorrentNotFoundInClientError) Handle() transport.HTTPResponse {
	return buildTorrentNotFoundResponse()
}

func buildTorrentNotFoundResponse() transport.HTTPResponse {
	return transport.NewHTTPResponse(http.StatusBadRequest, common.ApiError{
		Status:    http.StatusBadRequest,
		ErrorCode: string(TorrentNotFoundInClientErrorCode),
		Message:   "You cannot resume torrent that not added",
	})
}
