package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TorrentHandler struct{}

func NewTorrentHandler() *TorrentHandler {
	return &TorrentHandler{}
}

func (h *TorrentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	torrents := rg.Group("/torrents")
	{
		torrents.POST("", h.AddTorrent)
		torrents.GET("", h.GetTorrents)
		torrents.GET("/:infoHash", h.GetTorrent)
		torrents.GET("/:infoHash/peers", h.GetTorrentPeers)
		torrents.GET("/:infoHash/signatures", h.GetTorrentSignatures)
		torrents.POST("/:infoHash/pause", h.PauseTorrent)
		torrents.POST("/:infoHash/resume", h.ResumeTorrent)
		torrents.POST("/:infoHash/sign", h.SignTorrent)
	}
}

func (h *TorrentHandler) AddTorrent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "torrent added"})
}

func (h *TorrentHandler) GetTorrents(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *TorrentHandler) GetTorrent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *TorrentHandler) GetTorrentPeers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *TorrentHandler) GetTorrentSignatures(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *TorrentHandler) PauseTorrent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *TorrentHandler) ResumeTorrent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *TorrentHandler) SignTorrent(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
