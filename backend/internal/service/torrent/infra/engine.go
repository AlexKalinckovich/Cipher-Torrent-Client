package infra

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/peer"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/reputation"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/anacrolix/torrent/peer_protocol"
)

type AnacrolixEngine struct {
	client    *torrent.Client
	publisher ports.EventPublisher
	mapper    *MessageMapper
}

func NewAnacrolixEngine(dataDir string, publisher ports.EventPublisher) (*AnacrolixEngine, error) {
	engine := &AnacrolixEngine{publisher: publisher, mapper: NewMessageMapper()}
	return engine.initClient(dataDir)
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

func (e *AnacrolixEngine) StartDownload(ctx context.Context, infoBytes []byte) (string, error) {
	t, err := e.addTorrent(infoBytes)
	if err != nil {
		return "", err
	}
	return e.startAndReturnHash(t)
}

func (e *AnacrolixEngine) addTorrent(infoBytes []byte) (*torrent.Torrent, error) {
	mi := &metainfo.MetaInfo{InfoBytes: infoBytes}
	return e.client.AddTorrent(mi)
}

func (e *AnacrolixEngine) startAndReturnHash(t *torrent.Torrent) (string, error) {
	t.DownloadAll()
	return t.InfoHash().HexString(), nil
}

func (e *AnacrolixEngine) onReadMessage(pc *torrent.PeerConn, msg *peer_protocol.Message) {
	e.publishEvent(e.getChannel(pc), e.mapper.ToPacketEvent(pc, msg))
}

func (e *AnacrolixEngine) onExtendedHandshake(pc *torrent.PeerConn, hs *peer_protocol.ExtendedHandshakeMessage) {
	e.publishEvent(e.getChannel(pc), e.buildPeerEvent(pc, hs))
}

func (e *AnacrolixEngine) onExtensionMessage(evt torrent.PeerConnReadExtensionMessageEvent) {
	e.publishEvent(e.getChannel(evt.PeerConn), e.buildReputationEvent(evt))
}

func (e *AnacrolixEngine) getChannel(pc *torrent.PeerConn) string {
	return fmt.Sprintf("torrent_%s", pc.Torrent().InfoHash().HexString())
}

func (e *AnacrolixEngine) publishEvent(channel string, event any) {
	payload, err := json.Marshal(event)
	if err != nil {
		return
	}
	_ = e.publisher.Publish(context.Background(), channel, payload)
}

func (e *AnacrolixEngine) buildPeerEvent(pc *torrent.PeerConn, hs *peer_protocol.ExtendedHandshakeMessage) peer.ConnectedEvent {
	return peer.ConnectedEvent{Type: peer.EventTypePeerConnected, InfoHash: pc.Torrent().InfoHash().HexString(), Peer: e.buildPeerModel(pc, hs)}
}

func (e *AnacrolixEngine) buildPeerModel(pc *torrent.PeerConn, hs *peer_protocol.ExtendedHandshakeMessage) peer.Model {
	return peer.Model{
		Identity: peer.Identity{ID: pc.PeerID.String(), Client: hs.V},
		Flags:    peer.Flags{},
		State:    peer.State{},
	}
}

func (e *AnacrolixEngine) buildReputationEvent(evt torrent.PeerConnReadExtensionMessageEvent) reputation.UpdateEvent {
	return reputation.UpdateEvent{
		Type:               reputation.EventTypeReputationUpdate,
		From:               reputation.PeerID(evt.PeerConn.PeerID.String()),
		To:                 "",
		NewReputationScore: 0,
	}
}
