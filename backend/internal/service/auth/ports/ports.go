package auth

import (
	"context"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/service/user"
	"time"

	userModel "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/model/user"
)

type UserDomainPort interface {
	Create(ctx context.Context, params user.CreateUserInput) (userModel.UserFull, error)
	GetByEmail(ctx context.Context, email string) (userModel.UserFull, error)
	Get(ctx context.Context, id int64) (userModel.UserFull, error)
	GetPasswordHash(ctx context.Context, email string) (string, error)
}

type RefreshTokenStorePort interface {
	Save(ctx context.Context, token string, userID int64, ttl time.Duration) error
	GetUserID(ctx context.Context, token string) (int64, error)
	Delete(ctx context.Context, token string) error
}

type TokenIssuerPort interface {
	IssueAccessToken(userID int64, peerID string) (string, error)
	IssueRefreshToken() (string, error)
}

type PasswordHasherPort interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}
