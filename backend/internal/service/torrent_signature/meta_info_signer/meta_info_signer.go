package signing

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent_signature/ports"
	"github.com/anacrolix/torrent/bencode"
)

type MetaInfoSigner struct{}

func NewMetaInfoSigner() *MetaInfoSigner {
	return &MetaInfoSigner{}
}

func (s *MetaInfoSigner) InjectSignature(req ports.InjectionRequest) ([]byte, error) {
	sigDict := map[string]interface{}{
		"ed25519_pubkey": req.PubKey,
		"signature":      req.Signature,
		"timestamp":      req.Timestamp,
	}

	existing, ok := req.RootDict["signatures"].([]interface{})
	if !ok {
		existing = []interface{}{}
	}

	req.RootDict["signatures"] = append(existing, sigDict)

	return bencode.Marshal(req.RootDict)
}
