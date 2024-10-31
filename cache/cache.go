package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"sync"
)

type Cache interface {
	Set(ctx context.Context, key, response string) error
	Get(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type InMemoryCache struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		data: make(map[string]string),
	}
}

func (c *InMemoryCache) Set(ctx context.Context, key, response string) error {
	hashedKey := hashKey(key)
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[hashedKey] = response
	return nil
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (string, error) {
	hashedKey := hashKey(key)
	c.mu.RLock()
	defer c.mu.RUnlock()
	response, exists := c.data[hashedKey]
	if !exists {
		return "", nil
	}
	return response, nil
}

func (c *InMemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	hashedKey := hashKey(key)
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, exists := c.data[hashedKey]
	return exists, nil
}

func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}
