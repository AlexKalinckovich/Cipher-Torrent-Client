package torrent_signature

import "time"

type TorrentSignatureEntity struct {
	ID            int64
	TorrentHash   []byte
	UserID        int64
	SignatureBlob []byte
	PayloadHash   []byte
	CreatedAt     time.Time
}
