package websocket

import (
	"context"
	"log"
	"sync"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/ports"
)

type Client struct {
	ID      string
	Send    chan []byte
	Channel string
}

type Hub struct {
	broker        ports.EventBroker
	clients       map[string]map[string]*Client
	activeStreams map[string]context.CancelFunc
	mu            sync.RWMutex
}

func NewHub(broker ports.EventBroker) *Hub {
	return &Hub{
		broker:        broker,
		clients:       make(map[string]map[string]*Client),
		activeStreams: make(map[string]context.CancelFunc),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[c.Channel] == nil {
		h.clients[c.Channel] = make(map[string]*Client)
	}
	h.clients[c.Channel][c.ID] = c
	log.Printf("[HUB] Registered client %s to channel %s. Total clients: %d", c.ID, c.Channel, len(h.clients[c.Channel]))

	h.ensureStreamListening(c.Channel)
}

func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients[c.Channel], c.ID)
	log.Printf("[HUB] Unregistered client %s from channel %s. Remaining: %d", c.ID, c.Channel, len(h.clients[c.Channel]))

	if len(h.clients[c.Channel]) == 0 {
		delete(h.clients, c.Channel)
		h.stopStreamListening(c.Channel)
	}
}

func (h *Hub) ensureStreamListening(channel string) {
	if _, exists := h.activeStreams[channel]; exists {
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.activeStreams[channel] = cancel
	log.Printf("[HUB] Starting stream listener for channel: %s", channel)
	go h.consumeStream(ctx, channel)
}

func (h *Hub) stopStreamListening(channel string) {
	if cancel, exists := h.activeStreams[channel]; exists {
		cancel()
		delete(h.activeStreams, channel)
		log.Printf("[HUB] Stopped stream listener for channel: %s", channel)
	}
}

func (h *Hub) consumeStream(ctx context.Context, channel string) {
	log.Printf("[HUB] 📡 Calling broker.Subscribe for channel: %s", channel)

	events, err := h.broker.Subscribe(ctx, channel)
	if err != nil {
		log.Printf("[HUB] ❌ Error subscribing to channel %s: %v", channel, err)
		return
	}

	log.Printf("[HUB] ✅ Successfully subscribed to channel: %s. Waiting for events...", channel)

	for {
		select {
		case <-ctx.Done():
			log.Printf("[HUB] Stream listener for channel %s stopped", channel)
			return
		case evt, ok := <-events:
			if !ok {
				log.Printf("[HUB] ❌ Channel %s closed by broker", channel)
				return
			}
			log.Printf("[HUB] 📤 Broadcasting event to %d clients on %s", len(h.clients[evt.Channel]), evt.Channel)
			h.broadcast(evt)
		}
	}
}

func (h *Hub) broadcast(evt ports.Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := h.clients[evt.Channel]
	for _, client := range clients {
		select {
		case client.Send <- evt.Payload:
			// Successfully sent
		default:
			log.Printf("[HUB] Client %s send channel full, dropping message", client.ID)
		}
	}
}
