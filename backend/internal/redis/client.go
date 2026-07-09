package redis

import (
	"context"
	"time"

	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/redis_errors"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(addr string, password string) (*redis.Client, error) {
	client := buildRedisClient(addr, password)
	return client, verifyConnection(client)
}

func buildRedisClient(addr string, password string) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})
}

func verifyConnection(client *redis.Client) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := client.Ping(ctx).Err()
	if err != nil {
		return redis_errors.NewRedisFailedConnectionErrorWithCause("failed to connect to Redis", err)
	}
	return nil
}
