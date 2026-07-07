package redis

import (
	"context"
	ports "github.com/AlexKalinckovich/Cipher-Torrent-Client/backend/internal/redis/ports"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

const SubscribeChannelCapacity = 100

type EventBroker struct {
	client *redis.Client
}

func NewEventBroker(addr string, password string) (ports.EventBroker, error) {
	client := redis.NewClient(&redis.Options{
		Addr:        addr,
		Password:    password,
		PoolSize:    20,
		DialTimeout: 5 * time.Second,
		ReadTimeout: 3 * time.Second,
	})
	return &EventBroker{client: client}, client.Ping(context.Background()).Err()
}

func (b *EventBroker) Publish(ctx context.Context, channel string, payload []byte) error {
	err := b.client.XAdd(ctx, &redis.XAddArgs{
		Stream: channel,
		Values: map[string]any{"payload": payload},
	}).Err()

	if err != nil {
		log.Printf("Redis XAdd failed: %v", err)
	}

	return err
}

func (b *EventBroker) Subscribe(ctx context.Context, channel string) (<-chan ports.Event, error) {
	ch := make(chan ports.Event, SubscribeChannelCapacity)
	go b.streamReader(ctx, channel, ch)
	return ch, nil
}

func (b *EventBroker) streamReader(ctx context.Context, channel string, ch chan<- ports.Event) {
	lastID := "$"
	for {
		msgs, err := b.readStream(ctx, channel, lastID)
		if err != nil || ctx.Err() != nil {
			return
		}
		b.dispatchMessages(msgs, ch, &lastID)
	}
}

func (b *EventBroker) readStream(ctx context.Context, channel string, lastID string) ([]redis.XMessage, error) {
	res, err := b.client.XRead(ctx, &redis.XReadArgs{
		Streams: []string{channel, lastID}, Count: 10, Block: 5 * time.Second,
	}).Result()
	if err != nil || len(res) == 0 {
		return nil, err
	}
	return res[0].Messages, nil
}

func (b *EventBroker) dispatchMessages(msgs []redis.XMessage, ch chan<- ports.Event, lastID *string) {
	for _, msg := range msgs {
		if evt, ok := b.buildEvent(msg); ok {
			ch <- evt
			*lastID = msg.ID
		}
	}
}

func (b *EventBroker) buildEvent(msg redis.XMessage) (ports.Event, bool) {
	str, ok := msg.Values["payload"].(string)
	return ports.Event{ID: msg.ID, Payload: []byte(str)}, ok
}

func (b *EventBroker) mapToEvents(msgs []redis.XMessage) []ports.Event {
	events := make([]ports.Event, 0, len(msgs))
	for _, msg := range msgs {
		if evt, ok := b.buildEvent(msg); ok {
			events = append(events, evt)
		}
	}
	return events
}

func (b *EventBroker) Poll(ctx context.Context, channel string, lastID string) ([]ports.Event, error) {
	msgs, err := b.readStream(ctx, channel, lastID)
	return b.mapToEvents(msgs), err
}
