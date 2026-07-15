package ports

import "context"

type StorageUploadRequest struct {
	InfoHash      []byte
	CreatorPubKey []byte
	FileBytes     []byte
}

type StorageIdentityRequest struct {
	InfoHash      []byte
	CreatorPubKey []byte
}

type TorrentStoragePort interface {
	UploadBaseTorrent(ctx context.Context, req StorageUploadRequest) error
	DownloadBaseTorrent(ctx context.Context, req StorageIdentityRequest) ([]byte, error)
	DeleteTorrent(ctx context.Context, req StorageIdentityRequest) error
}
