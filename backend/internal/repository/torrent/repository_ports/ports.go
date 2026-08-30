package repository_ports

import (
	"context"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/models"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type TorrentIdentityRepositoryRequest = models.TorrentIdentity
type CreateTorrentRepositoryRequest struct {
	Entity        torrentModel.TorrentEntity
	CreatorPubKey []byte
	CreatorUserID int64
}

type CreateUserTorrentRepositoryRequest struct {
	UserID        int64
	InfoHash      []byte
	CreatorPubKey []byte
	Status        torrentModel.TorrentStatus
}

type UserTorrentIdentityRepositoryRequest = models.UserTorrentIdentity
type UpdateStatusRepositoryRequest = models.UpdateTorrentStatusRequest
type UpdateProgressRepositoryRequest = models.UpdateTorrentProgressRequest

type TorrentRepositoryPort interface {
	CreateTorrent(ctx context.Context, req CreateTorrentRepositoryRequest) error
	CreateUserTorrent(ctx context.Context, req CreateUserTorrentRepositoryRequest) error
	GetTorrentByIdentity(ctx context.Context, req TorrentIdentityRepositoryRequest) (torrentModel.TorrentEntity, error)
	GetUserTorrents(ctx context.Context, userID int64) ([]torrentModel.TorrentDTO, error)
	UpdateUserTorrentStatus(ctx context.Context, req UpdateStatusRepositoryRequest) error
	UpdateUserTorrentProgress(ctx context.Context, req UpdateProgressRepositoryRequest) error
	DeleteUserTorrent(ctx context.Context, req UserTorrentIdentityRepositoryRequest) error
}
