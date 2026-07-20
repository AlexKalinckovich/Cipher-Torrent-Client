package torrent

import (
	"time"
)

type EventType string

const (
	EventTypeProgress   EventType = "torrent:progress"
	EventTypeSigned     EventType = "torrent:signed"
	EventTypePeer       EventType = "peer:connected"
	EventTypeReputation EventType = "reputation:update"
	EventTypePacket     EventType = "packet"
)

type TorrentDTO struct {
	InfoHash         string         `json:"info_hash"`
	CreatorPublicKey string         `json:"creator_public_key"`
	Name             string         `json:"name"`
	SizeBytes        int64          `json:"size_bytes"`
	StoragePath      string         `json:"storage_path"`
	Status           TorrentStatus  `json:"status"`
	Progress         float32        `json:"progress"`
	AddedAt          time.Time      `json:"added_at"`
	Files            []FileDTO      `json:"files"`
	Signatures       []SignatureDTO `json:"signatures"`
}

type FileDTO struct {
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
}

type SignatureDTO struct {
	SignerPublicKey string `json:"signer_public_key"`
	SignatureBytes  string `json:"signature_bytes"`
	Timestamp       int64  `json:"timestamp"`
}

type AddTorrentRequest struct {
	MagnetURI string `json:"magnet_uri,omitempty"`
	SavePath  string `json:"save_path"`
}

type AddTorrentResponse struct {
	Status   string `json:"status"`
	InfoHash string `json:"info_hash"`
}

type ProgressEvent struct {
	Type             EventType `json:"event_type"`
	InfoHash         string    `json:"info_hash"`
	Progress         float32   `json:"progress"`
	DownloadSpeedBps int64     `json:"download_speed_bps"`
	UploadSpeedBps   int64     `json:"upload_speed_bps"`
}

type SignedEvent struct {
	Type      EventType    `json:"event_type"`
	InfoHash  string       `json:"info_hash"`
	Signature SignatureDTO `json:"signature"`
}
