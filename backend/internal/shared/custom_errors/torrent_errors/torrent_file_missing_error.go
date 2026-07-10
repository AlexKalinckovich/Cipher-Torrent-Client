package torrent_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const TorrentFileMissingErrorCode abstract_error_code.ErrorCode = "TORRENT_FILE_MISSING"

type TorrentFileMissingError struct{}

func NewTorrentFileMissingError() *TorrentFileMissingError {
	return &TorrentFileMissingError{}
}

func (e *TorrentFileMissingError) Error() string {
	return "torrent file is missing in the request"
}

func (e *TorrentFileMissingError) Code() abstract_error_code.ErrorCode {
	return TorrentFileMissingErrorCode
}

func (e *TorrentFileMissingError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateBadRequestResponse(TorrentFileMissingErrorCode, e.Error())
}
