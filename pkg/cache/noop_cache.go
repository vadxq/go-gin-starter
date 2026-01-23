package cache

import (
	"context"
	"time"
)

// NoopCache is a no-op cache implementation for graceful degradation.
type NoopCache struct{}

// NewNoop creates a no-op cache instance.
func NewNoop() Cache {
	return &NoopCache{}
}

// Get reads from cache.
func (n *NoopCache) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, ErrNotFound
}

// Set writes to cache.
func (n *NoopCache) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return nil
}

// Delete removes a key from cache.
func (n *NoopCache) Delete(ctx context.Context, key string) error {
	return nil
}

// Clear clears cache.
func (n *NoopCache) Clear(ctx context.Context) error {
	return nil
}

// GetObject reads and decodes a cached object.
func (n *NoopCache) GetObject(ctx context.Context, key string, value interface{}) error {
	return ErrNotFound
}

// SetObject encodes and writes a cached object.
func (n *NoopCache) SetObject(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return nil
}
