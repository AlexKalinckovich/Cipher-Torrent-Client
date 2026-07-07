package handlers

import (
	"context"
	"encoding/hex"
	"errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

const torrentFileField = "torrent_file"

const publicKeyHeader = "X-Public-Key"

type TorrentServicePort interface {
	Inspect(ctx context.Context, file multipart.File) (torrent.Model, error)
	Download(ctx context.Context, file multipart.File) (ports.DownloadResponse, error)
}

type TorrentHandler struct {
	service TorrentServicePort
}

func NewTorrentHandler(service TorrentServicePort) *TorrentHandler {
	return &TorrentHandler{service: service}
}

func (h *TorrentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	torrents := rg.Group("/torrents")
	{
		torrents.POST("/inspect", h.Inspect)
		torrents.POST("/download", h.Download)
	}
}

func (h *TorrentHandler) Inspect(c *gin.Context) {
	file, _, err := c.Request.FormFile(torrentFileField)
	if err != nil {
		h.fail(c, err)
		return
	}
	defer file.Close()
	res, err := h.service.Inspect(c.Request.Context(), file)
	h.respond(c, http.StatusOK, res, err)
}

func (h *TorrentHandler) Download(c *gin.Context) {
	file, _, err := c.Request.FormFile(torrentFileField)
	if err != nil {
		h.fail(c, err)
		return
	}
	defer file.Close()

	res, err := h.service.Download(c.Request.Context(), file)
	h.respond(c, http.StatusCreated, res, err)
}

func (h *TorrentHandler) extractPubKey(c *gin.Context) ([]byte, error) {
	pubKeyHex := c.GetHeader(publicKeyHeader)
	if pubKeyHex == "" {
		return nil, errors.New("missing X-Public-Key header")
	}
	return hex.DecodeString(pubKeyHex)
}

func (h *TorrentHandler) respond(c *gin.Context, status int, data any, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(status, data)
}

func (h *TorrentHandler) fail(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
