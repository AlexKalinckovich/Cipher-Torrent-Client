package ports

import (
	"context"
)

type TorrentStoragePort interface {
	UploadBaseTorrent(ctx context.Context, infoHash []byte, fileBytes []byte) error
	DownloadBaseTorrent(ctx context.Context, infoHash []byte) ([]byte, error)
	DeleteTorrent(ctx context.Context, infoHash []byte) error
}
