package token_store

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisClientPort interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
}

type RefreshTokenStore struct {
	client RedisClientPort
}

func NewRefreshTokenStore(client RedisClientPort) *RefreshTokenStore {
	return &RefreshTokenStore{client: client}
}

func (s *RefreshTokenStore) Save(ctx context.Context, token string, userID int64, ttl time.Duration) error {
	key := s.buildKey(token)
	return s.client.Set(ctx, key, userID, ttl).Err()
}

func (s *RefreshTokenStore) GetUserID(ctx context.Context, token string) (int64, error) {
	key := s.buildKey(token)
	return s.client.Get(ctx, key).Int64()
}

func (s *RefreshTokenStore) Delete(ctx context.Context, token string) error {
	key := s.buildKey(token)
	return s.client.Del(ctx, key).Err()
}

func (s *RefreshTokenStore) buildKey(token string) string {
	return "refresh_token:" + token
}
