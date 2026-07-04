package redis

import (
	"context"
	"fmt"
	"github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/ports"
	"time"

	"github.com/redis/go-redis/v9"
)

type EventBus struct {
	client *redis.Client
}

func NewEventBus(addr string, password string) (ports.EventPublisher, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		PoolSize:     20,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &EventBus{client: client}, nil
}

func (b *EventBus) Publish(
	ctx context.Context,
	channel string,
	payload []byte,
) error {
	err := b.client.Publish(ctx, channel, payload).Err()
	//fmt.Printf("Published a message to channel '%s' with payload '%s'\n err '%s'", channel, payload, err)
	return err
}
