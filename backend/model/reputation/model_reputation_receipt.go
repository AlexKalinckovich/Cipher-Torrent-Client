package reputation

import (
	"time"
)

type EventType string

const (
	EventTypeReputationUpdate EventType = "reputation:update"
)

type PeerID string
type Signature string

type Receipt struct {
	From       PeerID    `json:"from_peer_id"`
	To         PeerID    `json:"to_peer_id"`
	PieceIndex int32     `json:"piece_index"`
	ByteCount  int64     `json:"byte_count"`
	Signature  Signature `json:"signature"`
	Timestamp  time.Time `json:"timestamp"`
}

type UpdateEvent struct {
	Type               EventType `json:"event_type"`
	From               PeerID    `json:"from_peer_id"`
	To                 PeerID    `json:"to_peer_id"`
	NewReputationScore float32   `json:"new_reputation_score"`
}
