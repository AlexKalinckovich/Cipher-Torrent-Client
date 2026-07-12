package ports

import (
	"context"
	"crypto/ed25519"
	"mime/multipart"
)

type SignatureRequest struct {
	InfoHash     []byte
	AnnounceList [][]string
	PrivateKey   ed25519.PrivateKey
}

type SignatureResult struct {
	Signature   []byte
	PayloadHash []byte
	Timestamp   int64
}

type InjectionRequest struct {
	RootDict  map[string]interface{}
	PubKey    []byte
	Signature []byte
	Timestamp int64
}

type SignatureBuilder interface {
	BuildPayloadAndSign(req SignatureRequest) (*SignatureResult, error)
}

type MetaInfoSigner interface {
	InjectSignature(req InjectionRequest) ([]byte, error)
}

type TorrentSigningServicePort interface {
	SignTorrent(ctx context.Context, file multipart.File, userID int64) ([]byte, error)
}
