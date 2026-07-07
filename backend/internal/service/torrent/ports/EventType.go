package ports

import (
	packet "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/packet"
	reputation "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/reputation"
)

type PeerVerifiedEvent struct {
	PeerID string
	PubKey []byte
}

type ReputationReceiptEvent struct {
	Receipt reputation.Receipt
}

type PacketTelemetryEvent struct {
	Log packet.Log
}

type DownloadResponse struct {
	InfoHash string `json:"info_hash"`
}
