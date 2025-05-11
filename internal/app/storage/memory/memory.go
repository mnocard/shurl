package memoryStorage

import (
	"errors"
	"net/url"
	"sync"

	hash "github.com/mnocard/shurl/internal/app/hash"
)

type MemoryStorage struct {
	mu        sync.RWMutex
	addresses map[string]string
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		addresses: make(map[string]string),
	}
}

func (s *MemoryStorage) Get(hash string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	url, exists := s.addresses[hash]
	if !exists {
		return "", errors.New("url not found")
	}

	return url, nil
}

func (s *MemoryStorage) Add(u string) (string, error) {
	if u == "" {
		return "", errors.New("url is empty")
	}

	if _, err := url.ParseRequestURI(u); err != nil {
		return "", err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	h := hash.GetHash([]byte(u))
	s.addresses[h] = u
	return h, nil
}
