package storage

import (
	"context"
	"errors"
	"sync"
)

// InMemoryStorage is a simple in-memory map-based storage useful for tests.
type InMemoryStorage struct {
	mu sync.RWMutex
	m  map[string][]byte
}

// NewInMemoryStorage constructs a new in-memory storage.
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{m: make(map[string][]byte)}
}

func (s *InMemoryStorage) Put(ctx context.Context, key string, data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	b := make([]byte, len(data))
	copy(b, data)
	s.m[key] = b
	return nil
}

func (s *InMemoryStorage) Get(ctx context.Context, key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[key]
	if !ok {
		return nil, errors.New("not found")
	}
	b := make([]byte, len(v))
	copy(b, v)
	return b, nil
}

func (s *InMemoryStorage) Delete(ctx context.Context, key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, key)
	return nil
}

func (s *InMemoryStorage) Exists(ctx context.Context, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[key]
	return ok, nil
}
