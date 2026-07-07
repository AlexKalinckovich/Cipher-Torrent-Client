package infra

import "sync"

type PeerRegistry struct {
	keys sync.Map
}

func NewPeerRegistry() *PeerRegistry {
	return &PeerRegistry{}
}

func (r *PeerRegistry) SaveKey(peerID string, pubKey []byte) {
	r.keys.Store(peerID, pubKey)
}

func (r *PeerRegistry) GetKey(peerID string) []byte {
	if val, ok := r.keys.Load(peerID); ok {
		return val.([]byte)
	}
	return nil
}

func (r *PeerRegistry) RemoveKey(peerID string) {
	r.keys.Delete(peerID)
}
