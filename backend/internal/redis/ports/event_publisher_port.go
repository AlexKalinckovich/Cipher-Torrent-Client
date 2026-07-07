package ports

import "context"

type EventPublisher interface {
	Publish(ctx context.Context, channel string, payload []byte) error
}
