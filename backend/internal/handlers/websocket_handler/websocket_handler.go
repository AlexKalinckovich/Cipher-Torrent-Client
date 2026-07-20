package websocket_handler

import (
	"log"
	"net/http"

	websocketHub "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/infrastructure/websocket"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketHandler struct {
	hub *websocketHub.Hub
}

func NewWebSocketHandler(hub *websocketHub.Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/ws/torrent/:infoHash", h.Handle)
	rg.GET("/ws/dashboard", h.HandleDashboard)
}

func (h *WebSocketHandler) HandleDashboard(c *gin.Context) {
	log.Printf("[WS] Incoming Dashboard WebSocket connection request")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] Failed to upgrade dashboard connection: %v", err)
		return
	}

	client := &websocketHub.Client{
		ID:      generateID(),
		Send:    make(chan []byte, 256),
		Channel: "global_progress",
	}
	h.hub.Register(client)

	go h.writePump(conn, client)
	go h.readPump(conn, client)
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
	infoHash := c.Param("infoHash")
	log.Printf("[WS] Incoming WebSocket connection request for infoHash: %s", infoHash)

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] Failed to upgrade connection: %v", err)
		return
	}
	log.Printf("[WS] Connection upgraded successfully for infoHash: %s", infoHash)

	client := h.createClient(infoHash)
	h.hub.Register(client)

	go h.writePump(conn, client)
	go h.readPump(conn, client)
}

func (h *WebSocketHandler) createClient(infoHash string) *websocketHub.Client {
	return &websocketHub.Client{
		ID:      generateID(),
		Send:    make(chan []byte, 256),
		Channel: "torrent_" + infoHash,
	}
}

func generateID() string {
	return uuid.New().String()
}

func (h *WebSocketHandler) writePump(conn *websocket.Conn, client *websocketHub.Client) {
	defer conn.Close()
	log.Printf("[WS] Write pump started for client %s", client.ID)

	for msg := range client.Send {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("[WS] Write error for client %s: %v", client.ID, err)
			return
		}
	}
	log.Printf("[WS] Write pump stopped for client %s", client.ID)
}

func (h *WebSocketHandler) readPump(conn *websocket.Conn, client *websocketHub.Client) {
	defer h.hub.Unregister(client)
	defer close(client.Send)
	log.Printf("[WS] Read pump started for client %s", client.ID)

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			log.Printf("[WS] Read error for client %s: %v", client.ID, err)
			return
		}
	}
}
