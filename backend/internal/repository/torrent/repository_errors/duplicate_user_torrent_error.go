package repository_errors

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/default_error_handlers"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport"
)

const (
	DuplicateUserTorrentErrorCode abstract_error_code.ErrorCode = "USER_TORRENT_ALREADY_EXISTS"
)

type DuplicateUserTorrentError struct{}

func (e *DuplicateUserTorrentError) Code() abstract_error_code.ErrorCode {
	return DuplicateUserTorrentErrorCode
}

func (e *DuplicateUserTorrentError) Error() string {
	return "user already has this torrent"
}

func (e *DuplicateUserTorrentError) Handle() transport.HTTPResponse {
	return default_error_handlers.CreateConflictResponse(DuplicateUserTorrentErrorCode, e.Error())
}
