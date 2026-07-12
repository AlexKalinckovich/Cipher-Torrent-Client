package repository_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	DuplicateTorrentErrorCode abstract_error_code.ErrorCode = "TORRENT_ALREADY_EXISTS"
)

type DuplicateTorrentError struct{}

func (e *DuplicateTorrentError) Code() abstract_error_code.ErrorCode {
	return DuplicateTorrentErrorCode
}

func (e *DuplicateTorrentError) Error() string {
	return "torrent with this info_hash already exists"
}

func (e *DuplicateTorrentError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateConflictResponse(DuplicateTorrentErrorCode, e.Error())
}
