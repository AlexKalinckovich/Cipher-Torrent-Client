package infra

import "sync"

type PeerSessionStore struct {
	mu      sync.RWMutex
	pubKeys map[[20]byte][]byte
}

func NewPeerSessionStore() *PeerSessionStore {
	return &PeerSessionStore{pubKeys: make(map[[20]byte][]byte)}
}

func (s *PeerSessionStore) SaveKey(peerID [20]byte, pubKey []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pubKeys[peerID] = pubKey
}

func (s *PeerSessionStore) GetKey(peerID [20]byte) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key, ok := s.pubKeys[peerID]
	return key, ok
}
