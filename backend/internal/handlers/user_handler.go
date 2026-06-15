package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	// dependencies for working with MySQL 8.0
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) RegisterRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("/me", h.GetMyProfile)
		users.GET("/:peerID/reputation", h.GetPeerReputation)
	}
}

func (h *UserHandler) GetMyProfile(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}

func (h *UserHandler) GetPeerReputation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "OK"})
}
