package license

import "context"

type MemoryRepository struct {
	items map[string]License
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{items: make(map[string]License)}
}

func (s *MemoryRepository) Create(_ context.Context, l License) (License, error) {
	if l.Version <= 0 {
		l.Version = 1
	}
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
	cur, ok := s.items[id]
	if !ok {
		return License{}, ErrNotFound
	}
	if in.Version <= 0 {
		return License{}, ErrBadVersion
	}
	if in.Version != cur.Version {
		return License{}, ErrConflict
	}
	in.ID = id
	in.Version = cur.Version + 1
	s.items[id] = in
	return in, nil
}
