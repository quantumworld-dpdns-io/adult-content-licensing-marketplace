package license

import "context"

type MemoryRepository struct {
	items map[string]License
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{items: make(map[string]License)}
}

func (s *MemoryRepository) Create(_ context.Context, l License) (License, error) {
	s.items[l.ID] = l
	return l, nil
}

func (s *MemoryRepository) List(_ context.Context) ([]License, error) {
	out := make([]License, 0, len(s.items))
	for _, l := range s.items {
		out = append(out, l)
	}
	return out, nil
}

func (s *MemoryRepository) Get(_ context.Context, id string) (License, error) {
	l, ok := s.items[id]
	if !ok {
		return License{}, ErrNotFound
	}
	return l, nil
}

func (s *MemoryRepository) Update(_ context.Context, id string, in License) (License, error) {
	_, ok := s.items[id]
	if !ok {
		return License{}, ErrNotFound
	}
	in.ID = id
	s.items[id] = in
	return in, nil
}
