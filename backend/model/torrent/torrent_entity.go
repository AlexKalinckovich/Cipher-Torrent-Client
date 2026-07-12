package torrent

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/abstract_error_code"
	"time"
)

type TorrentStatus string

const (
	StatusIdle        TorrentStatus = "idle"
	StatusDownloading TorrentStatus = "downloading"
	StatusSeeding     TorrentStatus = "seeding"
	StatusPaused      TorrentStatus = "paused"
)

type TorrentEntity struct {
	InfoHash    []byte    `gorm:"primaryKey;type:binary(20)" json:"-"`
	InfoBytes   []byte    `gorm:"type:mediumblob" json:"-"`
	Name        string    `gorm:"type:varchar(255)" json:"name"`
	SizeBytes   int64     `json:"size_bytes"`
	PieceLength int       `json:"piece_length"`
	IsPrivate   bool      `json:"is_private"`
	StoragePath string    `gorm:"type:varchar(2048)" json:"storage_path"`
	AddedAt     time.Time `json:"added_at"`
}

const TorrentEntityCode abstract_error_code.ErrorCode = "TORRENT_NOT_FOUND"

func (e TorrentEntity) EntityCode() abstract_error_code.ErrorCode {
	return TorrentEntityCode
}

func (e TorrentEntity) EntityName() string {
	return "Torrent"
}

type UserTorrentEntity struct {
	UserID          int64         `gorm:"primaryKey" json:"user_id"`
	TorrentInfoHash []byte        `gorm:"primaryKey;type:binary(20)" json:"-"`
	Status          TorrentStatus `gorm:"type:varchar(20);default:idle" json:"status"`
	Progress        float32       `gorm:"type:decimal(5,2);default:0.00" json:"progress"`
}

const UserTorrentEntityCode abstract_error_code.ErrorCode = "USER_TORRENT_NOT_FOUND"

func (e UserTorrentEntity) EntityCode() abstract_error_code.ErrorCode {
	return UserTorrentEntityCode
}

func (e UserTorrentEntity) EntityName() string {
	return "UserTorrent"
}
