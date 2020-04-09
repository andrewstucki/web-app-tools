package state

import (
	"context"
	"errors"
	"sync"
)

// MemoryTokenManager just shoves the token
// some place in an in memory cache
type MemoryTokenManager struct {
	cache *sync.Map
}

// NewMemoryTokenManager initializes a MemoryTokenManager
func NewMemoryTokenManager() *MemoryTokenManager {
	return &MemoryTokenManager{
		cache: &sync.Map{},
	}
}

// Set stores the token in the cache
func (m *MemoryTokenManager) Set(ctx context.Context, subject, token string) error {
	m.cache.Store(subject, token)
	return nil
}

// Get returns a token from the cache
func (m *MemoryTokenManager) Get(ctx context.Context, subject string) (string, error) {
	value, ok := m.cache.Load(subject)
	if !ok {
		return "", errors.New("not found")
	}
	return value.(string), nil
}
