package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	RefreshTokenLength = 32
)

type JWTIssuer struct {
	secret []byte
	ttl    time.Duration
}

func NewJWTIssuer(secret []byte, ttl time.Duration) *JWTIssuer {
	return &JWTIssuer{secret: secret, ttl: ttl}
}

func (i *JWTIssuer) IssueAccessToken(userID int64, peerID string) (string, error) {
	claims := i.buildClaims(userID, peerID)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return i.signToken(token)
}

func (i *JWTIssuer) buildClaims(userID int64, peerID string) jwt.MapClaims {
	return jwt.MapClaims{
		"user_id": userID,
		"peer_id": peerID,
		"exp":     time.Now().Add(i.ttl).Unix(),
		"iat":     time.Now().Unix(),
	}
}

func (i *JWTIssuer) signToken(token *jwt.Token) (string, error) {
	return token.SignedString(i.secret)
}

func (i *JWTIssuer) IssueRefreshToken() (string, error) {
	bytes, err := i.generateRandomBytes()
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (i *JWTIssuer) generateRandomBytes() ([]byte, error) {
	b := make([]byte, RefreshTokenLength)
	_, err := rand.Read(b)
	return b, err
}
