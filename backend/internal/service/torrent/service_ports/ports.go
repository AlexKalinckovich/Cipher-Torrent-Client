package service_ports

import (
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
	"mime/multipart"
)

type CreateTorrentServiceRequest struct {
	File             multipart.File
	UserID           int64
	CreatorPublicKey []byte
}

type AddTorrentServiceRequest struct {
	InfoHash      []byte
	CreatorPubKey []byte
	UserID        int64
}

type TorrentIdentityServiceRequest struct {
	InfoHash      []byte
	CreatorPubKey []byte
	UserID        int64
}

type UpdateProgressServiceRequest struct {
	InfoHash      []byte
	CreatorPubKey []byte
	UserID        int64
	Progress      float32
}

type TorrentEngine interface {
	StartDownload(infoBytes []byte) (string, error)
	PauseTorrent(infoHash []byte) error
	ResumeTorrent(infoBytes []byte) error
}

type TorrentValidatorPort interface {
	ValidateAddRequest(req torrentModel.AddTorrentRequest) error
}
