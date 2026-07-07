package peer

import (
	"time"
)

type EventType string

const (
	EventTypePeerConnected EventType = "peer:connected"
)

type Identity struct {
	ID     string `json:"peer_id"`
	IP     string `json:"ip"`
	Port   int32  `json:"ports"`
	Client string `json:"client_name"`
}

type Flags struct {
	Encryption         bool `json:"encryption,omitempty"`
	UTP                bool `json:"utp,omitempty"`
	DHT                bool `json:"dht,omitempty"`
	PEX                bool `json:"pex,omitempty"`
	SupportsReputation bool `json:"supports_ut_reputation,omitempty"`
}

type State struct {
	ReputationScore  float32   `json:"reputation_score,omitempty"`
	LastSeen         time.Time `json:"last_seen,omitempty"`
	ConnectedSince   time.Time `json:"connected_since,omitempty"`
	DownloadSpeedBps int64     `json:"download_speed_bps,omitempty"`
	UploadSpeedBps   int64     `json:"upload_speed_bps,omitempty"`
}

type Model struct {
	Identity
	Flags
	State
}

type ConnectedEvent struct {
	Type     EventType `json:"event_type"`
	InfoHash string    `json:"info_hash"`
	Peer     Model     `json:"peer"`
}
