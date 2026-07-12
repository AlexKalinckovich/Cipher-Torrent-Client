package ports

import (
	"context"

	torrentSignatureModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent_signature"
)

type TorrentSignatureRepositoryPort interface {
	Create(ctx context.Context, entity torrentSignatureModel.TorrentSignatureEntity) error
	CreateInTransaction(ctx context.Context, entity torrentSignatureModel.TorrentSignatureEntity) error
	GetByTorrentHash(ctx context.Context, torrentHash []byte) ([]torrentSignatureModel.TorrentSignatureEntity, error)
	GetByUserAndTorrentHash(ctx context.Context, userID int64, torrentHash []byte) (*torrentSignatureModel.TorrentSignatureEntity, error)
	DeleteByTorrentHash(ctx context.Context, torrentHash []byte) error
}
