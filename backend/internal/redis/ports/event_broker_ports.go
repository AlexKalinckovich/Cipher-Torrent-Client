package ports

import "context"

type Event struct {
	ID      string
	Channel string
	Payload []byte
}

type EventBroker interface {
	Publish(ctx context.Context, channel string, payload []byte) error
	Subscribe(ctx context.Context, channel string) (<-chan Event, error)
	Poll(ctx context.Context, channel string, lastID string) ([]Event, error)
}
