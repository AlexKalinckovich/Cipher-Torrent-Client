package jwt

import (
	headererrors "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/shared/custom_errors/transport"
	"time"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

const (
	AccessTokenTTL = 15 * time.Minute
)

type JWTParser struct {
	secret []byte
}

func NewJWTParser(secret []byte) *JWTParser {
	return &JWTParser{secret: secret}
}

func (p *JWTParser) Parse(tokenString string) (*middleware.Claims, error) {
	token, err := p.parseToken(tokenString)
	if err != nil {
		return nil, headererrors.NewInvalidTokenError(err)
	}
	return p.extractClaims(token)
}

func (p *JWTParser) parseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, p.getKeyFunc())
}

func (p *JWTParser) getKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return p.secret, nil
	}
}

func (p *JWTParser) extractClaims(token *jwt.Token) (*middleware.Claims, error) {
	if !token.Valid {
		return nil, headererrors.NewInvalidTokenError(nil)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, headererrors.NewInvalidTokenClaimsError()
	}
	return p.mapToMiddlewareClaims(claims)
}

func (p *JWTParser) mapToMiddlewareClaims(claims jwt.MapClaims) (*middleware.Claims, error) {
	userID, err := p.extractUserID(claims)
	if err != nil {
		return nil, err
	}
	peerID, err := p.extractPeerID(claims)
	if err != nil {
		return nil, err
	}
	return &middleware.Claims{
		UserID: userID,
		PeerID: peerID,
	}, nil
}

func (p *JWTParser) extractUserID(claims jwt.MapClaims) (int64, error) {
	val, ok := claims["user_id"]
	if !ok {
		return 0, headererrors.NewInvalidTokenClaimsError()
	}
	floatVal, ok := val.(float64)
	if !ok {
		return 0, headererrors.NewInvalidTokenClaimsError()
	}
	return int64(floatVal), nil
}

func (p *JWTParser) extractPeerID(claims jwt.MapClaims) (string, error) {
	val, ok := claims["peer_id"]
	if !ok {
		return "", headererrors.NewInvalidTokenClaimsError()
	}
	strVal, ok := val.(string)
	if !ok {
		return "", headererrors.NewInvalidTokenClaimsError()
	}
	return strVal, nil
}
