package middleware

import (
	"sync"
	"time"
)

// MemoryStore 内存存储实现
type MemoryStore struct {
	mu    sync.RWMutex
	items map[string]memoryItem
}

type memoryItem struct {
	value   string
	expires time.Time
}

// NewMemoryStore 创建内存存储实例
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items: make(map[string]memoryItem),
	}
}

func (s *MemoryStore) Get(key string) (string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, found := s.items[key]
	if !found {
		return "", nil
	}

	if !item.expires.IsZero() && time.Now().After(item.expires) {
		delete(s.items, key)
		return "", nil
	}

	return item.value, nil
}

func (s *MemoryStore) Set(key, value string, expire time.Duration) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	var expires time.Time
	if expire > 0 {
		expires = time.Now().Add(expire)
	}

	s.items[key] = memoryItem{
		value:   value,
		expires: expires,
	}

	return nil
}

func (s *MemoryStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
	return nil
}

func (s *MemoryStore) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, found := s.items[key]
	return found
}
