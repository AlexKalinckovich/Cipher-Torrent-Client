package packet

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/reputation"
	"time"
)

type EventType string

const (
	EventTypePacket EventType = "packet"
)

type Direction string

const (
	DirectionInbound  Direction = "inbound"
	DirectionOutbound Direction = "outbound"
)

type MessageType string

const (
	MessageTypeChoke         MessageType = "choke"
	MessageTypeUnchoke       MessageType = "unchoke"
	MessageTypeInterested    MessageType = "interested"
	MessageTypeNotInterested MessageType = "not_interested"
	MessageTypeHandshake     MessageType = "handshake"
	MessageTypeBitfield      MessageType = "bitfield"
	MessageTypeHave          MessageType = "have"
	MessageTypeRequest       MessageType = "request"
	MessageTypePiece         MessageType = "piece"
	MessageTypeCancel        MessageType = "cancel"
	MessageTypeReputation    MessageType = "reputation"
)

type Log struct {
	ID          int64       `json:"id"`
	Timestamp   time.Time   `json:"timestamp"`
	Direction   Direction   `json:"direction"`
	MessageType MessageType `json:"message_type"`
	PeerIP      string      `json:"peer_ip,omitempty"`

	PieceIndex  *int32 `json:"piece_index,omitempty"`
	ChunkOffset *int32 `json:"chunk_offset,omitempty"`
	ChunkLength *int32 `json:"chunk_length,omitempty"`

	ParsedInfo        string              `json:"parsed_info,omitempty"`
	ReputationPayload *reputation.Receipt `json:"reputation_payload,omitempty"`
	SizeBytes         int64               `json:"size_bytes,omitempty"`
}

type Event struct {
	Type   EventType `json:"event_type"`
	Packet Log       `json:"packet"`
}
