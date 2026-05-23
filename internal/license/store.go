package license

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("license not found")

type License struct {
	ID                   string   `json:"id"`
	CreatorID            string   `json:"creator_id"`
	Title                string   `json:"title"`
	Description          string   `json:"description"`
	AITraingProhibited   bool     `json:"ai_training_prohibited"`
	BasePriceCents       int64    `json:"base_price_cents"`
	Currency             string   `json:"currency"`
	TerritoriesISO2Codes []string `json:"territories_iso2_codes"`
}

type Store struct {
	mu    sync.RWMutex
	items map[string]License
}

func NewStore() *Store {
	return &Store{items: make(map[string]License)}
}

func (s *Store) Create(l License) License {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[l.ID] = l
	return l
}

func (s *Store) List() []License {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]License, 0, len(s.items))
	for _, l := range s.items {
		out = append(out, l)
	}
	return out
}

func (s *Store) Get(id string) (License, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	l, ok := s.items[id]
	if !ok {
		return License{}, ErrNotFound
	}
	return l, nil
}
