package torrent_signature_handler

import (
	"fmt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/torrent_signature/torrent_signature_service_ports"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport/decoders/base64_decoder"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport/decoders/hex_decoder"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport/json_binder"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetSignedTorrentJSON struct {
	InfoHash      string `json:"info_hash"`
	CreatorPubKey string `json:"creator_pub_key"`
}

type signTorrentJSON struct {
	InfoHash      string `json:"info_hash"`
	CreatorPubKey string `json:"creator_pub_key"`
}

type TorrentSignatureHandler struct {
	service torrent_signature_service_ports.TorrentSigningServicePort
}

func NewTorrentSignatureHandler(service torrent_signature_service_ports.TorrentSigningServicePort) *TorrentSignatureHandler {
	return &TorrentSignatureHandler{service: service}
}

func (h *TorrentSignatureHandler) RegisterRoutes(rg *gin.RouterGroup) {
	torrents := rg.Group("/torrents")
	torrents.POST("/sign", h.Sign)
	torrents.GET("/download", h.GetSignedTorrent)
}

func (h *TorrentSignatureHandler) Sign(c *gin.Context) {
	req, err := h.buildSignRequest(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	h.executeSign(c, req)
}
func (h *TorrentSignatureHandler) GetSignedTorrent(c *gin.Context) {
	var jsonReq GetSignedTorrentJSON
	if err := c.ShouldBindJSON(&jsonReq); err != nil {
		h.fail(c, err)
		return
	}

	infoHash, err := hex_decoder.DecodeHexParam(jsonReq.InfoHash)
	if err != nil {
		h.fail(c, err)
		return
	}

	creatorPubKey, err := base64_decoder.DecodeBase64Param(jsonReq.CreatorPubKey)
	if err != nil {
		h.fail(c, err)
		return
	}

	req := torrent_signature_service_ports.GetSignedTorrentFileRequest{
		InfoHash:      infoHash,
		CreatorPubKey: creatorPubKey,
	}

	fileBytes, err := h.service.GetSignedTorrentFile(c.Request.Context(), req)
	if err != nil {
		h.fail(c, err)
		return
	}

	c.Header("Content-Type", "application/x-bittorrent")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.torrent"`, jsonReq.InfoHash))

	c.Data(http.StatusOK, "application/x-bittorrent", fileBytes)
}

func (h *TorrentSignatureHandler) buildSignRequest(c *gin.Context) (torrent_signature_service_ports.SignTorrentServiceRequest, error) {
	jsonReq := json_binder.BindJSON[signTorrentJSON](c)
	if jsonReq.Err != nil {
		return torrent_signature_service_ports.SignTorrentServiceRequest{}, jsonReq.Err
	}
	return h.mapToSignRequest(jsonReq.Value, c)
}

func (h *TorrentSignatureHandler) mapToSignRequest(jsonReq signTorrentJSON, c *gin.Context) (torrent_signature_service_ports.SignTorrentServiceRequest, error) {
	infoHash, err := hex_decoder.DecodeHexParam(jsonReq.InfoHash)
	if err != nil {
		return torrent_signature_service_ports.SignTorrentServiceRequest{}, err
	}
	return h.decodeCreatorAndBuild(infoHash, jsonReq.CreatorPubKey, c)
}

func (h *TorrentSignatureHandler) decodeCreatorAndBuild(infoHash []byte, creatorPubKeyHex string, c *gin.Context) (torrent_signature_service_ports.SignTorrentServiceRequest, error) {
	creatorPubKey, err := base64_decoder.DecodeBase64Param(creatorPubKeyHex)
	if err != nil {
		return torrent_signature_service_ports.SignTorrentServiceRequest{}, err
	}
	return h.buildFinalRequest(infoHash, creatorPubKey, c), nil
}

func (h *TorrentSignatureHandler) buildFinalRequest(infoHash []byte, creatorPubKey []byte, c *gin.Context) torrent_signature_service_ports.SignTorrentServiceRequest {
	return torrent_signature_service_ports.SignTorrentServiceRequest{
		InfoHash:      infoHash,
		CreatorPubKey: creatorPubKey,
		UserID:        h.extractUserID(c),
	}
}

func (h *TorrentSignatureHandler) executeSign(c *gin.Context, req torrent_signature_service_ports.SignTorrentServiceRequest) {
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
