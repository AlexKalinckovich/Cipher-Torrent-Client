package infra

import (
	"encoding/json"
	"fmt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/packet"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/reputation"
	"os"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/peer_protocol"
)

type MessageMapper struct {
	localPubKey []byte
}

func NewMessageMapper(localPubKey []byte) *MessageMapper {
	return &MessageMapper{
		localPubKey: localPubKey,
	}
}

func (m *MessageMapper) ToPacketEvent(pc *torrent.PeerConn, msg *peer_protocol.Message, key []byte) packet.Event {
	return packet.Event{Type: packet.EventTypePacket, Packet: m.buildLog(pc, msg, key)}
}

func (m *MessageMapper) buildLog(pc *torrent.PeerConn, msg *peer_protocol.Message, key []byte) packet.Log {
	log := m.initBaseLog(pc, msg)
	m.applyPieceMetadata(&log, msg)
	m.applyReputationPayload(&log, pc, msg, key)
	//m.writeLogToFile(log)
	return log
}

func (m *MessageMapper) writeLogToFile(log packet.Log) {
	jsonData, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		fmt.Printf("[ERROR] Failed to marshal log to JSON: %v\n", err)
		return
	}

	file, err := os.OpenFile("./local.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("[ERROR] Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()

	if _, err := file.Write(jsonData); err != nil {
		fmt.Printf("[ERROR] Failed to write log to file: %v\n", err)
		return
	}

	if _, err := file.Write([]byte(",\n")); err != nil {
		fmt.Printf("[ERROR] Failed to write separator: %v\n", err)
	}
}

func (m *MessageMapper) initBaseLog(pc *torrent.PeerConn, msg *peer_protocol.Message) packet.Log {
	return packet.Log{
		ID:          time.Now().UnixNano(),
		Timestamp:   time.Now(),
		Direction:   packet.DirectionInbound,
		MessageType: m.resolveType(msg.Type),
		PeerIP:      pc.RemoteAddr.String(),
		SizeBytes:   int64(msg.Length.Int()),
	}
}

func (m *MessageMapper) applyPieceMetadata(log *packet.Log, msg *peer_protocol.Message) {
	if !m.hasPieceInfo(msg) {
		return
	}
	idx := int32(msg.Index.Int())
	log.PieceIndex = &idx
	if m.isDataTransfer(msg) {
		begin := int32(msg.Begin.Int())
		length := int32(msg.Length.Int())
		log.ChunkOffset = &begin
		log.ChunkLength = &length
	}
}

func (m *MessageMapper) hasPieceInfo(msg *peer_protocol.Message) bool {
	return msg.Index != 0 || msg.Type == peer_protocol.Have || msg.Type == peer_protocol.Request
}

func (m *MessageMapper) isDataTransfer(msg *peer_protocol.Message) bool {
	return msg.Type == peer_protocol.Piece || msg.Type == peer_protocol.Request
}

func (m *MessageMapper) applyReputationPayload(log *packet.Log, pc *torrent.PeerConn, msg *peer_protocol.Message, key []byte) {
	if msg.Type != peer_protocol.Piece {
		return
	}
	fromPubKey := []byte(pc.PeerID.String())

	log.ReputationPayload = m.buildReceipt(pc, msg, fromPubKey, key)
}

func (m *MessageMapper) buildReceipt(pc *torrent.PeerConn, msg *peer_protocol.Message, fromPubKey []byte, key []byte) *reputation.Receipt {
	infoHash := pc.Torrent().InfoHash().HexString()
	pieceIdx := int32(msg.Index.Int())
	byteCount := int64(msg.Length.Int())

	return &reputation.Receipt{
		InfoHash:   infoHash,
		FromPubKey: fromPubKey,
		ToPubKey:   key,
		PieceIndex: pieceIdx,
		ByteCount:  byteCount,
		Timestamp:  time.Now().Unix(),
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
	return packet.MessageTypeReputation
}
