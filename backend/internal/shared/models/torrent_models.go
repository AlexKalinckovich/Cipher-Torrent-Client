package models

import (
	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
	"mime/multipart"
)

// TorrentIdentity represents the unique identity of a torrent
type TorrentIdentity struct {
	InfoHash      []byte
	CreatorPubKey []byte
}

// UserTorrentIdentity extends TorrentIdentity with user context
type UserTorrentIdentity struct {
	UserID        int64
	InfoHash      []byte
	CreatorPubKey []byte
}

// CreateTorrentRequest represents a request to create a new torrent
type CreateTorrentRequest struct {
	File             multipart.File
	UserID           int64
	CreatorPublicKey []byte
}

// AddTorrentRequest represents a request to add an existing torrent to user's list
type AddTorrentRequest struct {
	UserID        int64
	InfoHash      []byte
	CreatorPubKey []byte
}

// UpdateTorrentStatusRequest represents a request to update torrent status
type UpdateTorrentStatusRequest struct {
	UserID        int64
	InfoHash      []byte
	CreatorPubKey []byte
	Status        torrentModel.TorrentStatus
}

// UpdateTorrentProgressRequest represents a request to update torrent progress
type UpdateTorrentProgressRequest struct {
	UserID        int64
	InfoHash      []byte
	CreatorPubKey []byte
	Progress      float32
}

// ProgressUpdateBody represents the JSON body for progress updates
type ProgressUpdateBody struct {
	Progress float32 `json:"progress"`
}
