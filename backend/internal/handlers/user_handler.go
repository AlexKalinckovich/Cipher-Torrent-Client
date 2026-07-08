package handlers

import (
	"encoding/base64"
	serviceUser "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	userModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type bindResult[T any] struct {
	Value T
	Err   error
}

func bindJSON[T any](c *gin.Context) bindResult[T] {
	var value T
	err := c.ShouldBindJSON(&value)
	return bindResult[T]{Value: value, Err: err}
}

type UserHandler struct {
	service serviceUser.UserServicePort
}

func NewUserHandler(service serviceUser.UserServicePort) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.POST("", h.Create)
		users.GET("/:id", h.Get)
		users.GET("/email/:email", h.GetByEmail)
		users.GET("/nickname/:nickname", h.GetByNickname)
		users.GET("/public-key/*key", h.GetByPublicKey)
		users.PUT("/:id", h.Update)
		users.PATCH("/:id", h.Patch)
		users.DELETE("/:id", h.Delete)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	result := bindJSON[serviceUser.CreateUserInput](c)
	h.handleCreateBind(c, result)
}

func (h *UserHandler) handleCreateBind(c *gin.Context, result bindResult[serviceUser.CreateUserInput]) {
	if result.Err != nil {
		h.fail(c, result.Err)
		return
	}
	log.Printf(result.Value.Nickname)
	res, err := h.service.Create(c.Request.Context(), result.Value)
	h.handleCreateResult(c, res, err)
}

func (h *UserHandler) handleCreateResult(c *gin.Context, res userModel.UserFull, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	h.handleGetID(c, id, err)
}

func (h *UserHandler) handleGetID(c *gin.Context, id int64, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	res, serviceErr := h.service.Get(c.Request.Context(), id)
	h.handleGetResult(c, res, serviceErr)
}

func (h *UserHandler) handleGetResult(c *gin.Context, res userModel.UserFull, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) GetByEmail(c *gin.Context) {
	email := c.Param("email")
	res, err := h.service.GetByEmail(c.Request.Context(), email)
	h.handleGetResult(c, res, err)
}

func (h *UserHandler) GetByNickname(c *gin.Context) {
	nickname := c.Param("nickname")
	res, err := h.service.GetByNickname(c.Request.Context(), nickname)
	h.handleGetResult(c, res, err)
}

func (h *UserHandler) GetByPublicKey(c *gin.Context) {
	keyParam := c.Param("key")
	if len(keyParam) > 0 && keyParam[0] == '/' {
		keyParam = keyParam[1:]
	}

	keyBytes, err := base64.StdEncoding.DecodeString(keyParam)
	h.handlePublicKeyDecode(c, keyBytes, err)
}

func (h *UserHandler) handlePublicKeyDecode(c *gin.Context, keyBytes []byte, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	res, serviceErr := h.service.GetByPublicKey(c.Request.Context(), keyBytes)
	h.handleGetResult(c, res, serviceErr)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	h.handleUpdateID(c, id, err)
}

func (h *UserHandler) handleUpdateID(c *gin.Context, id int64, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	result := bindJSON[serviceUser.UpdateUserInput](c)
	h.handleUpdateBind(c, id, result)
}

func (h *UserHandler) handleUpdateBind(c *gin.Context, id int64, result bindResult[serviceUser.UpdateUserInput]) {
	if result.Err != nil {
		h.fail(c, result.Err)
		return
	}

	res, err := h.service.Update(c.Request.Context(), id, result.Value)
	h.handleUpdateResult(c, res, err)
}

func (h *UserHandler) handleUpdateResult(c *gin.Context, res userModel.UserFull, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Patch(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	h.handlePatchID(c, id, err)
}

func (h *UserHandler) handlePatchID(c *gin.Context, id int64, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	result := bindJSON[serviceUser.PatchUserFields](c)
	h.handlePatchBind(c, id, result)
}

func (h *UserHandler) handlePatchBind(c *gin.Context, id int64, result bindResult[serviceUser.PatchUserFields]) {
	if result.Err != nil {
		h.fail(c, result.Err)
		return
	}

	res, err := h.service.Patch(c.Request.Context(), id, result.Value)
	h.handlePatchResult(c, res, err)
}

func (h *UserHandler) handlePatchResult(c *gin.Context, res userModel.UserFull, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	h.handleDeleteID(c, id, err)
}

func (h *UserHandler) handleDeleteID(c *gin.Context, id int64, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	err = h.service.Delete(c.Request.Context(), id)
	h.handleDeleteResult(c, err)
}

func (h *UserHandler) handleDeleteResult(c *gin.Context, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) fail(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
