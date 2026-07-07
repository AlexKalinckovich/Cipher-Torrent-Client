package ports

import "context"

type TorrentEngine interface {
	StartDownload(ctx context.Context, infoBytes []byte) (string, error)
}
