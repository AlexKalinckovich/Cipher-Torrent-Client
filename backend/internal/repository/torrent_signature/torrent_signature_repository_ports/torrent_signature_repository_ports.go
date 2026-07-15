package torrent_signature_repository_ports

import (
	"context"
	torrentSignatureModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent_signature"
)

type TorrentIdentityRepositoryRequest struct {
	InfoHash      []byte
	CreatorPubKey []byte
}

type UserTorrentIdentityRepositoryRequest struct {
	UserID        int64
	InfoHash      []byte
	CreatorPubKey []byte
}

type CreateSignatureMapRequest struct {
	TorrentHash   []byte
	CreatorPubKey []byte
	SignerID      int64
	SignatureID   int64
	TrustLevel    int8
}

type TorrentSignatureRepositoryPort interface {
	CreateSignature(ctx context.Context, entity torrentSignatureModel.TorrentSignatureEntity) error
	CreateSignatureMap(ctx context.Context, req CreateSignatureMapRequest) error
	CreateInTransaction(ctx context.Context, sigEntity torrentSignatureModel.TorrentSignatureEntity, mapReq CreateSignatureMapRequest) error
	GetByTorrentIdentity(ctx context.Context, req TorrentIdentityRepositoryRequest) ([]torrentSignatureModel.TorrentSignatureEntity, error)
	GetByUserAndTorrentIdentity(ctx context.Context, req UserTorrentIdentityRepositoryRequest) (*torrentSignatureModel.TorrentSignatureEntity, error)
	DeleteByTorrentIdentity(ctx context.Context, req TorrentIdentityRepositoryRequest) error
}
