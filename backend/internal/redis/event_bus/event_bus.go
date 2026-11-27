package event_bus

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type EventBus struct {
	client *redis.Client
}

func (b *EventBus) Publish(ctx context.Context, channel string, payload []byte) error {
	return b.client.Publish(ctx, channel, payload).Err()
}
