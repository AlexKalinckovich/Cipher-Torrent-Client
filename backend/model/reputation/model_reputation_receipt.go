package reputation

type EventType string

const (
	EventTypeReputationUpdate EventType = "reputation:update"
)

type Receipt struct {
	InfoHash string `json:"info_hash" bencode:"info_hash"`

	FromPubKey []byte `json:"from_pub_key" bencode:"from_pub_key"`
	ToPubKey   []byte `json:"to_pub_key" bencode:"to_pub_key"`

	PieceIndex int32 `json:"piece_index" bencode:"piece_index"`
	ByteCount  int64 `json:"byte_count" bencode:"byte_count"`

	Timestamp int64 `json:"timestamp" bencode:"timestamp"`

	Signature []byte `json:"signature" bencode:"signature"`
}

type UpdateEvent struct {
	Type               EventType `json:"event_type"`
	FromPubKey         []byte    `json:"from_pub_key"`
	ToPubKey           []byte    `json:"to_pub_key"`
	NewReputationScore float32   `json:"new_reputation_score"`
}
