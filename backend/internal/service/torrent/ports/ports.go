package ports

import (
	"context"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type TorrentRepositoryPort interface {
	AddTorrent(ctx context.Context, entity torrentModel.TorrentEntity, userID int64) error
	GetByInfoHash(ctx context.Context, infoHash []byte) (torrentModel.TorrentEntity, error)
	GetUserTorrents(ctx context.Context, userID int64) ([]torrentModel.TorrentDTO, error)
	UpdateStatus(ctx context.Context, userID int64, infoHash []byte, status torrentModel.TorrentStatus) error
	UpdateProgress(ctx context.Context, userID int64, infoHash []byte, progress float32) error
	DeleteTorrent(ctx context.Context, infoHash []byte) error
}

type TorrentEngine interface {
	StartDownload(infoBytes []byte) (string, error)
	PauseTorrent(infoHash []byte) error
	ResumeTorrent(infoBytes []byte) error
}

type TorrentValidatorPort interface {
	ValidateAddRequest(req torrentModel.AddTorrentRequest) error
}
