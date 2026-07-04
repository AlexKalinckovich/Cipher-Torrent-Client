package packet

import (
	"time"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/reputation"
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
	ID                int64              `json:"id"`
	Timestamp         time.Time          `json:"timestamp"`
	Direction         Direction          `json:"direction"`
	MessageType       MessageType        `json:"message_type"`
	PeerIP            string             `json:"peer_ip,omitempty"`
	RawPayloadBase64  string             `json:"raw_payload_base64"`
	ParsedInfo        string             `json:"parsed_info,omitempty"`
	ReputationPayload reputation.Receipt `json:"reputation_payload,omitempty"`
	SizeBytes         int64              `json:"size_bytes,omitempty"`
}

type Event struct {
	Type   EventType `json:"event_type"`
	Packet Log       `json:"packet"`
}
