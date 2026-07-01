package torrent

import (
	"time"
)

type Status string

const (
	StatusPending  Status = "pending"
	StatusRunning  Status = "running"
	StatusFinished Status = "finished"
	StatusError    Status = "error"
)

type EventType string

const (
	EventTypeProgress   EventType = "torrent:progress"
	EventTypeSigned     EventType = "torrent:signed"
	EventTypePeer       EventType = "peer:connected"
	EventTypeReputation EventType = "reputation:update"
	EventTypePacket     EventType = "packet"
)

type Model struct {
	ID               int64     `json:"id"`
	InfoHash         string    `json:"info_hash"`
	Name             string    `json:"name"`
	SizeBytes        int64     `json:"size_bytes"`
	Status           Status    `json:"status"`
	Progress         float32   `json:"progress,omitempty"`
	DownloadSpeedBps int64     `json:"download_speed_bps,omitempty"`
	UploadSpeedBps   int64     `json:"upload_speed_bps,omitempty"`
	AddedAt          time.Time `json:"added_at"`
	PieceLength      int64     `json:"piece_length,omitempty"`
	Files            []File    `json:"files,omitempty"`
	SignaturesCount  int32     `json:"signatures_count,omitempty"`
	IsSignedByMe     bool      `json:"is_signed_by_me,omitempty"`
	PeersCount       int32     `json:"peers_count,omitempty"`
}

type AddRequest struct {
	MagnetURI  string `json:"magnet_uri,omitempty"`
	FileBase64 string `json:"file_base64,omitempty"`
	SavePath   string `json:"save_path,omitempty"`
}

type AddResponse struct {
	Status   string `json:"status"`
	InfoHash string `json:"info_hash"`
}

type File struct {
	Path      string `json:"path"`
	SizeBytes int64  `json:"size_bytes"`
}

type ProgressEvent struct {
	Type             EventType `json:"event_type"`
	InfoHash         string    `json:"info_hash"`
	Progress         float32   `json:"progress"`
	DownloadSpeedBps int64     `json:"download_speed_bps"`
	UploadSpeedBps   int64     `json:"upload_speed_bps"`
}

type Signature struct {
	SignerUserID    int64     `json:"signer_user_id"`
	SignerPublicKey string    `json:"signer_public_key"`
	SignatureBytes  string    `json:"signature_bytes"`
	IsValid         bool      `json:"is_valid"`
	SignedAt        time.Time `json:"signed_at"`
}

type SignedEvent struct {
	Type      EventType `json:"event_type"`
	InfoHash  string    `json:"info_hash"`
	Signature Signature `json:"signature"`
}
