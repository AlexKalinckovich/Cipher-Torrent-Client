package service_ports

import (
	"context"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/models"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"mime/multipart"
)

type CreateTorrentServiceRequest = models.CreateTorrentRequest
type AddTorrentServiceRequest = models.AddTorrentRequest
type TorrentIdentityServiceRequest = models.UserTorrentIdentity
type UpdateProgressServiceRequest = models.UpdateTorrentProgressRequest

type TorrentServicePort interface {
	Inspect(file multipart.File) (torrentModel.TorrentEntity, error)
	Create(ctx context.Context, req CreateTorrentServiceRequest) (torrentModel.TorrentDTO, error)
	Add(ctx context.Context, req AddTorrentServiceRequest) error
	GetByInfoHash(ctx context.Context, req TorrentIdentityServiceRequest) (torrentModel.TorrentEntity, error)
	GetUserTorrents(ctx context.Context, userID int64) ([]torrentModel.TorrentDTO, error)
	GetStoreTorrents(ctx context.Context) ([]torrentModel.StoreTorrentDTO, error)
	PauseTorrent(ctx context.Context, req TorrentIdentityServiceRequest) error
	ResumeTorrent(ctx context.Context, req TorrentIdentityServiceRequest) error
	UpdateProgress(ctx context.Context, req UpdateProgressServiceRequest) error
	DeleteTorrent(ctx context.Context, req TorrentIdentityServiceRequest) error
	DeletePublishedTorrent(ctx context.Context, req TorrentIdentityServiceRequest) error
}

type TorrentEngine interface {
	StartDownload(mi *metainfo.MetaInfo, key []byte) (string, error)
	PauseTorrent(infoHash []byte) error
	ResumeTorrent(mi *metainfo.MetaInfo, key []byte) error
}

type TorrentValidatorPort interface {
	ValidateAddRequest(req torrentModel.AddTorrentRequest) error
}
