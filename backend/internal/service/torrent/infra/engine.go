package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/torrent_errors"
	"log"
	"time"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"
	pp "github.com/anacrolix/torrent/peer_protocol"

	eventPorts "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/packet"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/reputation"
)

const dpkiExtensionName pp.ExtensionName = "ut_dpki"

type UTExtensionPayload struct {
	MsgType int                 `bencode:"msg_type"`
	PubKey  []byte              `bencode:"pub_key,omitempty"`
	Receipt *reputation.Receipt `bencode:"receipt,omitempty"`
}

type asyncEvent struct {
	Channel string
	Payload any
}

type AnacrolixEngine struct {
	client       *torrent.Client
	publisher    eventPorts.EventPublisher
	mapper       *MessageMapper
	sessionStore *PeerSessionStore
	localPubKey  []byte
	eventBus     chan asyncEvent
}

func NewAnacrolixEngine(dataDir string, localPubKey []byte, publisher eventPorts.EventPublisher) (*AnacrolixEngine, error) {
	engine := &AnacrolixEngine{
		publisher:    publisher,
		mapper:       NewMessageMapper(localPubKey),
		sessionStore: NewPeerSessionStore(),
		localPubKey:  localPubKey,
		eventBus:     make(chan asyncEvent, 5000),
	}
	go engine.dispatchEvents()
	return engine.initClient(dataDir)
}

func (e *AnacrolixEngine) StartDownload(infoBytes []byte) (string, error) {
	t, err := e.addTorrent(infoBytes)
	if err != nil {
		return "", err
	}
	return e.activateTorrent(t)
}

func (e *AnacrolixEngine) PauseTorrent(infoHash []byte) error {
	hash := e.toMetainfoHash(infoHash)
	t, ok := e.client.Torrent(hash)
	if !ok {
		return torrent_errors.NewTorrentNotFoundInClientError()
	}
	t.Drop()
	return nil
}

func (e *AnacrolixEngine) ResumeTorrent(infoBytes []byte) error {
	t, err := e.addTorrent(infoBytes)
	if err != nil {
		return err
	}
	_, err = e.activateTorrent(t)
	return err
}

func (e *AnacrolixEngine) addTorrent(infoBytes []byte) (*torrent.Torrent, error) {
	mi := &metainfo.MetaInfo{InfoBytes: infoBytes}
	return e.client.AddTorrent(mi)
}

func (e *AnacrolixEngine) activateTorrent(t *torrent.Torrent) (string, error) {
	<-t.GotInfo()
	t.DownloadAll()
	return t.InfoHash().HexString(), nil
}

func (e *AnacrolixEngine) toMetainfoHash(infoHash []byte) metainfo.Hash {
	var hash metainfo.Hash
	copy(hash[:], infoHash)
	return hash
}

func (e *AnacrolixEngine) initClient(dataDir string) (*AnacrolixEngine, error) {
	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = dataDir
	e.bindCallbacks(cfg)
	return e.createClient(cfg)
}

func (e *AnacrolixEngine) bindCallbacks(cfg *torrent.ClientConfig) {
	cfg.Callbacks.ReadMessage = e.onReadMessage
	cfg.Callbacks.ReadExtendedHandshake = e.onExtendedHandshake
	cfg.Callbacks.PeerConnAdded = append(cfg.Callbacks.PeerConnAdded, e.onPeerConnAdded)
	cfg.Callbacks.SentRequest = append(cfg.Callbacks.SentRequest, e.onSentRequest)
	cfg.Callbacks.PeerConnReadExtensionMessage = append(cfg.Callbacks.PeerConnReadExtensionMessage, e.onExtensionMessage)
}

func (e *AnacrolixEngine) createClient(cfg *torrent.ClientConfig) (*AnacrolixEngine, error) {
	client, err := torrent.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	e.client = client
	return e, nil
}

func (e *AnacrolixEngine) dispatchEvents() {
	for evt := range e.eventBus {
		e.publishToRedis(evt)
	}
}

func (e *AnacrolixEngine) publishToRedis(evt asyncEvent) {
	data, err := json.Marshal(evt.Payload)
	if err != nil {
		log.Printf("Marshal error: %v", err)
		return
	}
	_ = e.publisher.Publish(context.Background(), evt.Channel, data)
}

func (e *AnacrolixEngine) enqueueEvent(channel string, payload any) {
	if channel == "" {
		return
	}
	select {
	case e.eventBus <- asyncEvent{Channel: channel, Payload: payload}:
	default:
	}
}

func (e *AnacrolixEngine) onPeerConnAdded(pc *torrent.PeerConn) {
	pc.LocalLtepProtocolMap.AddUserProtocol(dpkiExtensionName)
}

func (e *AnacrolixEngine) onExtendedHandshake(pc *torrent.PeerConn, hs *pp.ExtendedHandshakeMessage) {
	if _, supported := hs.M[dpkiExtensionName]; supported {
		e.sendPubKey(pc)
	}
}

func (e *AnacrolixEngine) sendPubKey(pc *torrent.PeerConn) {
	payload := e.buildDPKIPayload()
	e.writeExtendedAndLog(pc, payload)
	e.enqueueOutboundHandshake(pc, payload)
}

func (e *AnacrolixEngine) buildDPKIPayload() []byte {
	msg := UTExtensionPayload{MsgType: 0, PubKey: e.localPubKey}
	data, _ := bencode.Marshal(msg)
	return data
}

func (e *AnacrolixEngine) writeExtendedAndLog(pc *torrent.PeerConn, payload []byte) {
	err := pc.WriteExtendedMessage(dpkiExtensionName, payload)
	if err != nil {
		log.Printf("WriteExtendedMessage error: %v", err)
	}
}

func (e *AnacrolixEngine) enqueueOutboundHandshake(pc *torrent.PeerConn, payload []byte) {
	logEvent := e.buildHandshakeLog(pc, payload)
	e.enqueueEvent(e.getChannel(pc), packet.Event{Type: packet.EventTypePacket, Packet: logEvent})
}

func (e *AnacrolixEngine) buildHandshakeLog(pc *torrent.PeerConn, payload []byte) packet.Log {
	return packet.Log{
		ID:          time.Now().UnixNano(),
		Timestamp:   time.Now(),
		Direction:   packet.DirectionOutbound,
		MessageType: packet.MessageTypeHandshake,
		PeerIP:      pc.RemoteAddr.String(),
		SizeBytes:   int64(6 + len(payload)),
		ParsedInfo:  "Sent DPKI Public Key",
	}
}

func (e *AnacrolixEngine) onExtensionMessage(evt torrent.PeerConnReadExtensionMessageEvent) {
	extName, _, err := evt.PeerConn.LocalLtepProtocolMap.LookupId(evt.ExtensionNumber)
	if err != nil || extName != dpkiExtensionName {
		return
	}
	e.processDPKIMessage(evt)
}

func (e *AnacrolixEngine) processDPKIMessage(evt torrent.PeerConnReadExtensionMessageEvent) {
	extPayload, err := e.unmarshalDPKIPayload(evt.Payload)
	if err != nil {
		return
	}
	e.saveKeyIfHandshake(evt.PeerConn.PeerID, extPayload)
}

func (e *AnacrolixEngine) unmarshalDPKIPayload(payload []byte) (UTExtensionPayload, error) {
	var extPayload UTExtensionPayload
	err := bencode.Unmarshal(payload, &extPayload)
	return extPayload, err
}

func (e *AnacrolixEngine) saveKeyIfHandshake(peerID [20]byte, payload UTExtensionPayload) {
	if payload.MsgType == 0 {
		e.sessionStore.SaveKey(peerID, payload.PubKey)
	}
}

func (e *AnacrolixEngine) onSentRequest(evt torrent.PeerRequestEvent) {
	logEvent := e.buildRequestLog(evt)
	e.setRequestMetadata(&logEvent, evt.Request)
	e.enqueueEvent(e.getPeerChannel(evt.Peer), packet.Event{
		Type:   packet.EventTypePacket,
		Packet: logEvent,
	})
}

func (e *AnacrolixEngine) buildRequestLog(evt torrent.PeerRequestEvent) packet.Log {
	return packet.Log{
		ID:          time.Now().UnixNano(),
		Timestamp:   time.Now(),
		Direction:   packet.DirectionOutbound,
		MessageType: packet.MessageTypeRequest,
		PeerIP:      evt.Peer.RemoteAddr.String(),
		SizeBytes:   int64(evt.Request.Length),
	}
}

func (e *AnacrolixEngine) setRequestMetadata(log *packet.Log, req torrent.Request) {
	pIdx, cOff, cLen := int32(req.Index), int32(req.Begin), int32(req.Length)
	log.PieceIndex, log.ChunkOffset, log.ChunkLength = &pIdx, &cOff, &cLen
}

func (e *AnacrolixEngine) onReadMessage(pc *torrent.PeerConn, msg *pp.Message) {
	pubKey, _ := e.sessionStore.GetKey(pc.PeerID)
	logEvent := e.mapper.ToPacketEvent(pc, msg, pubKey)
	e.enqueueEvent(e.getChannel(pc), logEvent)
}

func (e *AnacrolixEngine) getChannel(pc *torrent.PeerConn) string {
	if pc.Torrent() == nil {
		return ""
	}
	return fmt.Sprintf("torrent_%s", pc.Torrent().InfoHash().HexString())
}

func (e *AnacrolixEngine) getPeerChannel(p *torrent.Peer) string {
	if p.Torrent() == nil {
		return ""
	}
	return fmt.Sprintf("torrent_%s", p.Torrent().InfoHash().HexString())
}
