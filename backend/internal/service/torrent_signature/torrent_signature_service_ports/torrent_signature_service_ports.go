package torrent_signature_service_ports

import (
	"context"
	"crypto/ed25519"
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type SignTorrentServiceRequest struct {
	InfoHash      []byte
	CreatorPubKey []byte
	UserID        int64
}

type SignatureServiceRequest struct {
	InfoHash     []byte
	AnnounceList [][]string
	PrivateKey   ed25519.PrivateKey
}

type SignatureServiceResult struct {
	Signature   []byte
	PayloadHash []byte
	Timestamp   int64
}

type InjectionServiceRequest struct {
	RootDict  map[string]interface{}
	PubKey    []byte
	Signature []byte
	Timestamp int64
}

type GetSignedTorrentFileRequest struct {
	InfoHash      []byte
	CreatorPubKey []byte
}

type SignatureBuilder interface {
	BuildPayloadAndSign(req SignatureServiceRequest) (*SignatureServiceResult, error)
}

type MetaInfoSigner interface {
	InjectSignature(req InjectionServiceRequest) ([]byte, error)
}

type TorrentSigningServicePort interface {
	SignTorrent(ctx context.Context, req SignTorrentServiceRequest) (torrentModel.TorrentDTO, error)
	GetSignedTorrentFile(ctx context.Context, req GetSignedTorrentFileRequest) ([]byte, error)
}
