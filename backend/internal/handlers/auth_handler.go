package handlers

import (
	"context"
	"net/http"

	modelAuth "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/auth"
	"github.com/gin-gonic/gin"
)

type AuthServicePort interface {
	Login(ctx context.Context, email, password string) (modelAuth.AuthResponse, error)
	Register(ctx context.Context, email, password string) (modelAuth.AuthResponse, error)
	Refresh(ctx context.Context, refreshToken string) (modelAuth.TokensResponse, error)
}

type AuthHandler struct {
	service AuthServicePort
}

func NewAuthHandler(service AuthServicePort) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	authGroup := rg.Group("/auth")
	{
		authGroup.POST("/login", h.LoginUser)
		authGroup.POST("/register", h.RegisterUser)
		authGroup.POST("/refresh", h.RefreshToken)
	}
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	req, err := h.bindLoginRequest(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	result, serviceErr := h.service.Login(c.Request.Context(), req.Email, req.Password)
	h.respondAuth(c, result, serviceErr)
}

func (h *AuthHandler) bindLoginRequest(c *gin.Context) (modelAuth.AuthLoginRequest, error) {
	var req modelAuth.AuthLoginRequest
	err := c.ShouldBindJSON(&req)
	return req, err
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	req, err := h.bindRegisterRequest(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	result, serviceErr := h.service.Register(c.Request.Context(), req.Email, req.Password)
	h.respondAuth(c, result, serviceErr)
}

func (h *AuthHandler) bindRegisterRequest(c *gin.Context) (modelAuth.AuthRegisterRequest, error) {
	var req modelAuth.AuthRegisterRequest
	err := c.ShouldBindJSON(&req)
	return req, err
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	req, err := h.bindRefreshRequest(c)
	if err != nil {
		h.fail(c, err)
		return
	}
	result, serviceErr := h.service.Refresh(c.Request.Context(), req.RefreshToken)
	h.respondTokens(c, result, serviceErr)
}

func (h *AuthHandler) bindRefreshRequest(c *gin.Context) (modelAuth.RefreshTokenRequest, error) {
	var req modelAuth.RefreshTokenRequest
	err := c.ShouldBindJSON(&req)
	return req, err
}

func (h *AuthHandler) respondAuth(c *gin.Context, result modelAuth.AuthResponse, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) respondTokens(c *gin.Context, result modelAuth.TokensResponse, err error) {
	if err != nil {
		h.fail(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *AuthHandler) fail(c *gin.Context, err error) {
	_ = c.Error(err)
	c.Abort()
}
