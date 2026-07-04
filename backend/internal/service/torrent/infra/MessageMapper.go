package infra

import (
	"fmt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/packet"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/reputation"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/peer_protocol"
)

type MessageMapper struct{}

func NewMessageMapper() *MessageMapper {
	return &MessageMapper{}
}

func (m *MessageMapper) ToPacketEvent(pc *torrent.PeerConn, msg *peer_protocol.Message) packet.Event {
	return packet.Event{
		Type:   packet.EventTypePacket,
		Packet: m.buildLog(pc, msg),
	}
}

func (m *MessageMapper) buildLog(pc *torrent.PeerConn, msg *peer_protocol.Message) packet.Log {
	return packet.Log{
		ID:                time.Now().UnixNano(),
		Timestamp:         time.Now(),
		Direction:         packet.DirectionInbound,
		MessageType:       m.resolveType(msg.Type),
		PeerIP:            pc.RemoteAddr.String(),
		RawPayloadBase64:  "",
		ParsedInfo:        "",
		ReputationPayload: reputation.Receipt{},
		SizeBytes:         int64(msg.Length.Int()),
	}
}

func (m *MessageMapper) resolveType(t peer_protocol.MessageType) packet.MessageType {
	if mapped, ok := m.baseTypes()[t]; ok {
		return mapped
	}
	return m.resolveExtended(t)
}

func (m *MessageMapper) baseTypes() map[peer_protocol.MessageType]packet.MessageType {
	return map[peer_protocol.MessageType]packet.MessageType{
		peer_protocol.Choke:         packet.MessageTypeChoke,
		peer_protocol.Unchoke:       packet.MessageTypeUnchoke,
		peer_protocol.Interested:    packet.MessageTypeInterested,
		peer_protocol.NotInterested: packet.MessageTypeNotInterested,
		peer_protocol.Have:          packet.MessageTypeHave,
		peer_protocol.Bitfield:      packet.MessageTypeBitfield,
		peer_protocol.Request:       packet.MessageTypeRequest,
		peer_protocol.Piece:         packet.MessageTypePiece,
		peer_protocol.Cancel:        packet.MessageTypeCancel,
	}
}

func (m *MessageMapper) resolveExtended(t peer_protocol.MessageType) packet.MessageType {
	if t == peer_protocol.Extended {
		return packet.MessageTypeHandshake
	}
	return packet.MessageType(fmt.Sprintf("unknown_%d", t))
}
