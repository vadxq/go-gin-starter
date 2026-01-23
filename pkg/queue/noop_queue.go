package queue

import (
	"context"
	"time"
)

// NoopQueue is a no-op queue implementation for graceful degradation.
type NoopQueue struct{}

// NewNoop creates a no-op queue instance.
func NewNoop() Queue {
	return &NoopQueue{}
}

// Publish publishes a message.
func (n *NoopQueue) Publish(ctx context.Context, topic string, payload interface{}) error {
	return nil
}

// Subscribe registers a handler.
func (n *NoopQueue) Subscribe(ctx context.Context, topic string, handler Handler) error {
	return nil
}

// PublishDelayed publishes a delayed message.
func (n *NoopQueue) PublishDelayed(ctx context.Context, topic string, payload interface{}, delay time.Duration) error {
	return nil
}

// Close closes the queue.
func (n *NoopQueue) Close() error {
	return nil
}
