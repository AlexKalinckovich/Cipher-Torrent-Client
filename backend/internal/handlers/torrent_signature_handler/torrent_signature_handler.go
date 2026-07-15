package torrent_signature_handler

import (
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent_signature/ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport/decoders/base64_decoder"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport/decoders/hex_decoder"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport/json_binder"
	"net/http"

	"github.com/gin-gonic/gin"
)

type signTorrentJSON struct {
	InfoHash      string `json:"info_hash"`
	CreatorPubKey string `json:"creator_pub_key"`
}

type TorrentSignatureHandler struct {
	service ports.TorrentSigningServicePort
}

func NewTorrentSignatureHandler(service ports.TorrentSigningServicePort) *TorrentSignatureHandler {
	return &TorrentSignatureHandler{service: service}
}

func (h *TorrentSignatureHandler) RegisterRoutes(rg *gin.RouterGroup) {
	torrents := rg.Group("/torrents")
	torrents.POST("/sign", h.Sign)
}

func (h *TorrentSignatureHandler) Sign(c *gin.Context) {
	req, err := h.buildSignRequest(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	h.executeSign(c, req)
}

func (h *TorrentSignatureHandler) buildSignRequest(c *gin.Context) (ports.SignTorrentServiceRequest, error) {
	jsonReq := json_binder.BindJSON[signTorrentJSON](c)
	if jsonReq.Err != nil {
		return ports.SignTorrentServiceRequest{}, jsonReq.Err
	}
	return h.mapToSignRequest(jsonReq.Value, c)
}

func (h *TorrentSignatureHandler) mapToSignRequest(jsonReq signTorrentJSON, c *gin.Context) (ports.SignTorrentServiceRequest, error) {
	infoHash, err := hex_decoder.DecodeHexParam(jsonReq.InfoHash)
	if err != nil {
		return ports.SignTorrentServiceRequest{}, err
	}
	return h.decodeCreatorAndBuild(infoHash, jsonReq.CreatorPubKey, c)
}

func (h *TorrentSignatureHandler) decodeCreatorAndBuild(infoHash []byte, creatorPubKeyHex string, c *gin.Context) (ports.SignTorrentServiceRequest, error) {
	creatorPubKey, err := base64_decoder.DecodeBase64Param(creatorPubKeyHex)
	if err != nil {
		return ports.SignTorrentServiceRequest{}, err
	}
	return h.buildFinalRequest(infoHash, creatorPubKey, c), nil
}

func (h *TorrentSignatureHandler) buildFinalRequest(infoHash []byte, creatorPubKey []byte, c *gin.Context) ports.SignTorrentServiceRequest {
	return ports.SignTorrentServiceRequest{
		InfoHash:      infoHash,
		CreatorPubKey: creatorPubKey,
		UserID:        h.extractUserID(c),
	}
}

func (h *TorrentSignatureHandler) executeSign(c *gin.Context, req ports.SignTorrentServiceRequest) {
	dto, err := h.service.SignTorrent(c.Request.Context(), req)
	h.respond(c, http.StatusOK, dto, err)
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
