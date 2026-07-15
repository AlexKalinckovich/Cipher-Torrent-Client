package torrent_handler

import (
	"context"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent/service_ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport/decoders/base64_decoder"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport/decoders/hex_decoder"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport/json_binder"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"

	torrentModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/torrent"
)

const (
	torrentFileField   = "torrent_file"
	infoHashParam      = "info_hash"
	creatorPubKeyParam = "creator_pub_key"
	userIDKey          = "user_id"
	peerIDKey          = "peer_id"
)

type progressRequest struct {
	Progress float32 `json:"progress"`
}

type AddTorrentHandlerRequest struct {
	InfoHash      string `json:"info_hash"`
	CreatorPubKey string `json:"creator_pub_key"`
}

type TorrentServicePort interface {
	Inspect(file multipart.File) (torrentModel.TorrentEntity, error)
	Create(ctx context.Context, req service_ports.CreateTorrentServiceRequest) (torrentModel.TorrentDTO, error)
	Add(ctx context.Context, req service_ports.AddTorrentServiceRequest) error
	GetByInfoHash(ctx context.Context, req service_ports.TorrentIdentityServiceRequest) (torrentModel.TorrentEntity, error)
	GetUserTorrents(ctx context.Context, userID int64) ([]torrentModel.TorrentDTO, error)
	PauseTorrent(ctx context.Context, req service_ports.TorrentIdentityServiceRequest) error
	ResumeTorrent(ctx context.Context, req service_ports.TorrentIdentityServiceRequest) error
	UpdateProgress(ctx context.Context, req service_ports.UpdateProgressServiceRequest) error
	DeleteTorrent(ctx context.Context, req service_ports.TorrentIdentityServiceRequest) error
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
	torrents.POST("/create", h.Create)
	torrents.POST("/add", h.Add)
	torrents.GET("/", h.GetUserTorrents)
	torrents.GET("/:"+infoHashParam+"/:"+creatorPubKeyParam, h.GetByInfoHash)
	torrents.POST("/:"+infoHashParam+"/:"+creatorPubKeyParam+"/pause", h.PauseTorrent)
	torrents.POST("/:"+infoHashParam+"/:"+creatorPubKeyParam+"/resume", h.ResumeTorrent)
	torrents.POST("/:"+infoHashParam+"/:"+creatorPubKeyParam+"/progress", h.UpdateProgress)
	torrents.DELETE("/", h.DeleteTorrent)
}

func (h *TorrentHandler) Create(c *gin.Context) {
	file, err := h.extractFile(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	defer file.Close()
	h.buildAndExecuteCreate(c, file)

}

func (h *TorrentHandler) buildAndExecuteCreate(c *gin.Context, file multipart.File) {
	creatorPubKey, err := h.extractCreatorPubKey(c)
	if err != nil {
		log.Printf("Error: %v", err.Error())
		h.fail(c, err)
		return
	}
	req := service_ports.CreateTorrentServiceRequest{
		UserID:           h.extractUserID(c),
		CreatorPublicKey: creatorPubKey,
		File:             file,
	}
	res, err := h.service.Create(c.Request.Context(), req)
	h.respond(c, http.StatusCreated, res, err)
}

func (h *TorrentHandler) Add(c *gin.Context) {
	identity, err := h.extractIdentity(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	h.executeAdd(c, identity)
}

func (h *TorrentHandler) executeAdd(c *gin.Context, identity service_ports.TorrentIdentityServiceRequest) {
	req := service_ports.AddTorrentServiceRequest{
		InfoHash:      identity.InfoHash,
		CreatorPubKey: identity.CreatorPubKey,
		UserID:        h.extractUserID(c),
	}
	err := h.service.Add(c.Request.Context(), req)
	h.respondEmpty(c, http.StatusOK, err)
}

func (h *TorrentHandler) GetUserTorrents(c *gin.Context) {
	userID := h.extractUserID(c)
	res, err := h.service.GetUserTorrents(c.Request.Context(), userID)
	h.respond(c, http.StatusOK, res, err)
}

func (h *TorrentHandler) GetByInfoHash(c *gin.Context) {
	identity, err := h.extractIdentity(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	res, err := h.service.GetByInfoHash(c.Request.Context(), identity)
	h.respond(c, http.StatusOK, res, err)
}

func (h *TorrentHandler) PauseTorrent(c *gin.Context) {
	identity, err := h.extractIdentity(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	err = h.service.PauseTorrent(c.Request.Context(), identity)
	h.respondEmpty(c, http.StatusOK, err)
}

func (h *TorrentHandler) ResumeTorrent(c *gin.Context) {
	identity, err := h.extractIdentity(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	err = h.service.ResumeTorrent(c.Request.Context(), identity)
	h.respondEmpty(c, http.StatusOK, err)
}

func (h *TorrentHandler) UpdateProgress(c *gin.Context) {
	identity, err := h.extractIdentity(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	h.processProgressUpdate(c, identity)
}

func (h *TorrentHandler) processProgressUpdate(c *gin.Context, identity service_ports.TorrentIdentityServiceRequest) {
	var reqBody progressRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		h.fail(c, err)
		return
	}
	req := service_ports.UpdateProgressServiceRequest{
		UserID:        h.extractUserID(c),
		InfoHash:      identity.InfoHash,
		CreatorPubKey: identity.CreatorPubKey,
		Progress:      reqBody.Progress,
	}
	err := h.service.UpdateProgress(c.Request.Context(), req)
	h.respondEmpty(c, http.StatusOK, err)
}

func (h *TorrentHandler) DeleteTorrent(c *gin.Context) {
	userIdentity, err := h.extractUserTorrentIdentity(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	err = h.service.DeleteTorrent(c.Request.Context(), userIdentity)
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

func (h *TorrentHandler) extractCreatorPubKey(c *gin.Context) ([]byte, error) {
	val, _ := c.Get(peerIDKey)
	hexStr, ok := val.(string)
	log.Printf(hexStr)
	if !ok {
		return nil, nil
	}
	return []byte(hexStr), nil
}

func (h *TorrentHandler) extractIdentity(c *gin.Context) (service_ports.TorrentIdentityServiceRequest, error) {
	bindRes := json_binder.BindJSON[AddTorrentHandlerRequest](c)
	if bindRes.Err != nil {
		return service_ports.TorrentIdentityServiceRequest{}, bindRes.Err
	}

	infoHash, err := hex_decoder.DecodeHexParam(bindRes.Value.InfoHash)
	if err != nil {
		return service_ports.TorrentIdentityServiceRequest{}, err
	}

	creatorPubKey, err := base64_decoder.DecodeBase64Param(bindRes.Value.CreatorPubKey)
	if err != nil {
		return service_ports.TorrentIdentityServiceRequest{}, err
	}

	return service_ports.TorrentIdentityServiceRequest{
		InfoHash:      infoHash,
		CreatorPubKey: creatorPubKey,
	}, nil
}

func (h *TorrentHandler) extractUserTorrentIdentity(c *gin.Context) (service_ports.TorrentIdentityServiceRequest, error) {
	identity, err := h.extractIdentity(c)
	if err != nil {
		return service_ports.TorrentIdentityServiceRequest{}, err
	}
	return service_ports.TorrentIdentityServiceRequest{
		UserID:        h.extractUserID(c),
		InfoHash:      identity.InfoHash,
		CreatorPubKey: identity.CreatorPubKey,
	}, nil
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
