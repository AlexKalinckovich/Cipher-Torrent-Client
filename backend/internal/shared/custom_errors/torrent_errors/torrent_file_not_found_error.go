package torrent_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const TorrentFileNotFoundErrorCode abstract_error_code.ErrorCode = "TORRENT_FILE_NOT_FOUND"

type TorrentFileNotFoundError struct {
	Path string
}

func NewTorrentFileNotFoundError(path string) *TorrentFileNotFoundError {
	return &TorrentFileNotFoundError{Path: path}
}

func (e *TorrentFileNotFoundError) Error() string {
	return "file not found in torrent: " + e.Path
}

func (e *TorrentFileNotFoundError) Code() abstract_error_code.ErrorCode {
	return TorrentFileNotFoundErrorCode
}

func (e *TorrentFileNotFoundError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateBadRequestResponse(TorrentFileNotFoundErrorCode, e.Error())
}