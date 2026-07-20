package user_handler

import (
	"encoding/base64"
	serviceUser "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/user_errors"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/transport/json_binder"
	userModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/user"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

const (
	userIDKey = "user_id"
)

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
		users.GET("/me", h.GetCurrentUser)
		users.GET("/email/:email", h.GetByEmail)
		users.GET("/nickname/:nickname", h.GetByNickname)
		users.GET("/password/:email", h.GetPasswordHash)
		users.GET("/public-key", h.GetByPublicKey)
		users.PUT("/:id", h.Update)
		users.PATCH("/:id", h.Patch)
		users.DELETE("/:id", h.Delete)
	}
}

func (h *UserHandler) Create(c *gin.Context) {
	result := json_binder.BindJSON[serviceUser.CreateUserInput](c)
	h.handleCreateBind(c, result)
}

func (h *UserHandler) handleCreateBind(c *gin.Context, result json_binder.BindResult[serviceUser.CreateUserInput]) {
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

func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userId := h.extractUserID(c)
	res, err := h.service.Get(c.Request.Context(), userId)
	h.handleGetResult(c, res, err)
}

func (h *UserHandler) extractUserID(c *gin.Context) int64 {
	val, _ := c.Get(userIDKey)
	return val.(int64)
}

func (h *UserHandler) handleGetID(c *gin.Context, id int64, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	res, serviceErr := h.service.Get(c.Request.Context(), id)
	h.handleGetResult(c, res, serviceErr)
}

func (h *UserHandler) handleGetResult(c *gin.Context, res any, err error) {
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

func (h *UserHandler) GetPasswordHash(c *gin.Context) {
	email := c.Param("email")
	res, err := h.service.GetPasswordHash(c.Request.Context(), email)
	h.handleGetResult(c, res, err)
}

func (h *UserHandler) GetByPublicKey(c *gin.Context) {
	keyParam := c.Query("pub_key")

	if keyParam == "" {
		h.fail(c, user_errors.NewPublicKeyNotProvidedError())
		return
	}

	keyBytes, err := base64.RawURLEncoding.DecodeString(keyParam)

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
	result := json_binder.BindJSON[serviceUser.UpdateUserInput](c)
	h.handleUpdateBind(c, id, result)
}

func (h *UserHandler) handleUpdateBind(c *gin.Context, id int64, result json_binder.BindResult[serviceUser.UpdateUserInput]) {
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
	result := json_binder.BindJSON[serviceUser.PatchUserFields](c)
	h.handlePatchBind(c, id, result)
}

func (h *UserHandler) handlePatchBind(c *gin.Context, id int64, result json_binder.BindResult[serviceUser.PatchUserFields]) {
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
