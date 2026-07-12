package middleware

import (
	headererrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/token_errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type Claims struct {
	UserID int64
	PeerID string
}

type TokenParser interface {
	Parse(tokenString string) (*Claims, error)
}

type AuthMiddleware struct {
	parser TokenParser
}

func NewAuthMiddleware(parser TokenParser) *AuthMiddleware {
	return &AuthMiddleware{parser: parser}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		m.authenticate(c)
	}
}

func (m *AuthMiddleware) authenticate(c *gin.Context) {
	tokenString, err := m.extractToken(c)
	if err != nil {
		m.abortWithError(c, err)
		return
	}
	m.processToken(c, tokenString)
}

func (m *AuthMiddleware) processToken(c *gin.Context, tokenString string) {
	claims, err := m.parser.Parse(tokenString)
	if err != nil {
		m.abortWithError(c, err)
		return
	}
	m.setContext(c, claims)
	c.Next()
}

func (m *AuthMiddleware) extractToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", headererrors.NewMissingAuthHeaderError()
	}
	return m.parseBearerToken(authHeader)
}

func (m *AuthMiddleware) parseBearerToken(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", headererrors.NewInvalidAuthFormatError()
	}
	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

func (m *AuthMiddleware) abortWithError(c *gin.Context, err error) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
}

func (m *AuthMiddleware) setContext(c *gin.Context, claims *Claims) {
	c.Set("user_id", claims.UserID)
	c.Set("peer_id", claims.PeerID)
}
