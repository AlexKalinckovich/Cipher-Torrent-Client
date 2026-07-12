package torrent_signature_handler

import (
	"context"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

type TorrentSigningServicePort interface {
	SignTorrent(ctx context.Context, file multipart.File, userID int64) (torrentModel.TorrentDTO, error)
}

type TorrentSignatureHandler struct {
	service TorrentSigningServicePort
}

func NewTorrentSignatureHandler(service TorrentSigningServicePort) *TorrentSignatureHandler {
	return &TorrentSignatureHandler{service: service}
}

func (h *TorrentSignatureHandler) RegisterRoutes(rg *gin.RouterGroup) {
	torrents := rg.Group("/torrents")
	torrents.POST("/sign", h.Sign)
}

func (h *TorrentSignatureHandler) Sign(c *gin.Context) {
	file, err := h.extractFile(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	defer file.Close()
	h.processSigning(c, file)
}

func (h *TorrentSignatureHandler) processSigning(c *gin.Context, file multipart.File) {
	userID := h.extractUserID(c)
	dto, err := h.service.SignTorrent(c.Request.Context(), file, userID)
	h.respond(c, http.StatusOK, dto, err)
}

func (h *TorrentSignatureHandler) extractFile(c *gin.Context) (multipart.File, error) {
	file, _, err := c.Request.FormFile("torrent_file")
	return file, err
}

func (h *TorrentSignatureHandler) extractUserID(c *gin.Context) int64 {
	val, _ := c.Get("user_id")
	return val.(int64)
}

func (h *TorrentSignatureHandler) respond(c *gin.Context, status int, data any, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(status, data)
}

func (h *TorrentSignatureHandler) fail(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
