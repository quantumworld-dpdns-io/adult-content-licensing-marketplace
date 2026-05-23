package license

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	items map[string]License
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{items: make(map[string]License)}
}

func (s *MemoryRepository) Create(_ context.Context, l License) (License, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[l.ID] = l
	return l, nil
}

func (s *MemoryRepository) List(_ context.Context) ([]License, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]License, 0, len(s.items))
	for _, l := range s.items {
		out = append(out, l)
	}
	return out, nil
}

func (s *MemoryRepository) Get(_ context.Context, id string) (License, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.items[id]
	if !ok {
		return License{}, ErrNotFound
	}
	return l, nil
}
