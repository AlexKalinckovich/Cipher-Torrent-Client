package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type AuthHandler struct{}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

func (h *AuthHandler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/auth")
	{
		auth.POST("/login", h.LoginUser)
		auth.POST("/register", h.RegisterUser)
		auth.POST("/refresh", h.RefreshToken)
	}
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
