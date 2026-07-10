package torrent_handler

import (
	"context"
	"encoding/hex"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/torrent_errors"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

const (
	torrentFileField = "torrent_file"
	savePathField    = "save_path"
	infoHashParam    = "info_hash"
	userIDKey        = "user_id"
)

type progressRequest struct {
	Progress float32 `json:"progress"`
}

type TorrentServicePort interface {
	Inspect(file multipart.File) (torrentModel.TorrentEntity, error)
	Add(ctx context.Context, file multipart.File, userID int64, savePath string) (torrentModel.TorrentDTO, error)
	GetByInfoHash(ctx context.Context, infoHash []byte) (torrentModel.TorrentEntity, error)
	GetUserTorrents(ctx context.Context, userID int64) ([]torrentModel.TorrentDTO, error)
	PauseTorrent(ctx context.Context, userID int64, infoHash []byte) error
	ResumeTorrent(ctx context.Context, userID int64, infoHash []byte) error
	UpdateProgress(ctx context.Context, userID int64, infoHash []byte, progress float32) error
	DeleteTorrent(ctx context.Context, userID int64, infoHash []byte) error
}

type TorrentHandler struct {
	service      TorrentServicePort
	retranslator *TorrentHandlerErrorRetranslator
}

func NewTorrentHandler(service TorrentServicePort) *TorrentHandler {
	return &TorrentHandler{
		service:      service,
		retranslator: NewTorrentHandlerErrorRetranslator(),
	}
}

func (h *TorrentHandler) RegisterRoutes(rg *gin.RouterGroup) {
	torrents := rg.Group("/torrents")
	torrents.POST("/inspect", h.Inspect)
	torrents.POST("/add", h.Add)
	torrents.GET("/", h.GetUserTorrents)
	torrents.GET("/:"+infoHashParam, h.GetByInfoHash)
	torrents.POST("/:"+infoHashParam+"/pause", h.PauseTorrent)
	torrents.POST("/:"+infoHashParam+"/resume", h.ResumeTorrent)
	torrents.POST("/:"+infoHashParam+"/progress", h.UpdateProgress)
	torrents.DELETE("/:"+infoHashParam+"/delete", h.DeleteTorrent)
}

func (h *TorrentHandler) Inspect(c *gin.Context) {
	file, err := h.extractFile(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	defer file.Close()
	res, err := h.service.Inspect(file)
	h.respond(c, http.StatusOK, res, err)
}

func (h *TorrentHandler) Add(c *gin.Context) {
	file, err := h.extractFile(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	defer file.Close()
	userID := h.extractUserID(c)
	savePath := c.PostForm(savePathField)
	if savePath == "" {
		h.fail(c, torrent_errors.NewSafePathNotSpecifiedError())
	}
	res, err := h.service.Add(c.Request.Context(), file, userID, savePath)
	h.respond(c, http.StatusCreated, res, err)
}

func (h *TorrentHandler) GetUserTorrents(c *gin.Context) {
	userID := h.extractUserID(c)
	res, err := h.service.GetUserTorrents(c.Request.Context(), userID)
	h.respond(c, http.StatusOK, res, err)
}

func (h *TorrentHandler) GetByInfoHash(c *gin.Context) {
	infoHash, err := h.extractInfoHash(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	res, err := h.service.GetByInfoHash(c.Request.Context(), infoHash)
	h.respond(c, http.StatusOK, res, err)
}

func (h *TorrentHandler) PauseTorrent(c *gin.Context) {
	infoHash, userID, err := h.extractPauseResumeParams(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	err = h.service.PauseTorrent(c.Request.Context(), userID, infoHash)
	h.respondEmpty(c, http.StatusOK, err)
}

func (h *TorrentHandler) ResumeTorrent(c *gin.Context) {
	infoHash, userID, err := h.extractPauseResumeParams(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	err = h.service.ResumeTorrent(c.Request.Context(), userID, infoHash)
	h.respondEmpty(c, http.StatusOK, err)
}

func (h *TorrentHandler) UpdateProgress(c *gin.Context) {
	infoHash, userID, progress, err := h.extractProgressParams(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	err = h.service.UpdateProgress(c.Request.Context(), userID, infoHash, progress)
	h.respondEmpty(c, http.StatusOK, err)
}

func (h *TorrentHandler) extractFile(c *gin.Context) (multipart.File, error) {
	file, _, err := c.Request.FormFile(torrentFileField)
	if err != nil {
		return nil, h.retranslator.Retranslate(err)
	}
	return file, nil
}

func (h *TorrentHandler) extractUserID(c *gin.Context) int64 {
	val, _ := c.Get(userIDKey)
	return val.(int64)
}

func (h *TorrentHandler) extractInfoHash(c *gin.Context) ([]byte, error) {
	hexStr := c.Param(infoHashParam)
	return hex.DecodeString(hexStr)
}

func (h *TorrentHandler) extractProgress(c *gin.Context) (float32, error) {
	var req progressRequest
	err := c.ShouldBindJSON(&req)
	return req.Progress, err
}

func (h *TorrentHandler) extractPauseResumeParams(c *gin.Context) ([]byte, int64, error) {
	infoHash, err := h.extractInfoHash(c)
	if err != nil {
		return nil, 0, err
	}
	userID := h.extractUserID(c)
	return infoHash, userID, nil
}

func (h *TorrentHandler) extractProgressParams(c *gin.Context) ([]byte, int64, float32, error) {
	infoHash, err := h.extractInfoHash(c)
	if err != nil {
		return nil, 0, 0, err
	}
	userID := h.extractUserID(c)
	progress, err := h.extractProgress(c)
	return infoHash, userID, progress, err
}

func (h *TorrentHandler) respond(c *gin.Context, status int, data any, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(status, data)
}

func (h *TorrentHandler) respondEmpty(c *gin.Context, status int, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.Status(status)
}

func (h *TorrentHandler) fail(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}

func (h *TorrentHandler) DeleteTorrent(c *gin.Context) {
	infoHash, userID, err := h.extractPauseResumeParams(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	err = h.service.DeleteTorrent(c.Request.Context(), userID, infoHash)
	h.respondEmpty(c, http.StatusOK, err)
}
