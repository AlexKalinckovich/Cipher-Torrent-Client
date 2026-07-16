package signing

import (
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/binary"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent_signature/torrent_signature_service_ports"
	"time"

	"github.com/anacrolix/torrent/bencode"
)

type SignatureBuilder struct{}

func NewSignatureBuilder() *SignatureBuilder {
	return &SignatureBuilder{}
}

func (b *SignatureBuilder) BuildPayloadAndSign(req torrent_signature_service_ports.SignatureServiceRequest) (*torrent_signature_service_ports.SignatureServiceResult, error) {
	announceHash, err := b.hashAnnounceList(req.AnnounceList)
	if err != nil {
		return nil, err
	}

	timestamp := time.Now().Unix()
	timestampBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timestampBytes, uint64(timestamp))

	payload := make([]byte, 0, len(req.InfoHash)+len(announceHash)+len(timestampBytes))
	payload = append(payload, req.InfoHash...)
	payload = append(payload, announceHash...)
	payload = append(payload, timestampBytes...)

	payloadHash := sha256.Sum256(payload)
	signature := ed25519.Sign(req.PrivateKey, payload)

	return &torrent_signature_service_ports.SignatureServiceResult{
		Signature:   signature,
		PayloadHash: payloadHash[:],
		Timestamp:   timestamp,
	}, nil
}

func (b *SignatureBuilder) hashAnnounceList(announceList [][]string) ([]byte, error) {
	bencoded, err := bencode.Marshal(announceList)
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(bencoded)
	return hash[:], nil
}
