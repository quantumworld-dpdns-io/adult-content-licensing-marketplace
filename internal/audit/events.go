package audit

import (
	"sync"
	"time"
)

type Event struct {
	Type      string         `json:"type"`
	ActorSub  string         `json:"actor_sub"`
	ActorRole string         `json:"actor_role"`
	EntityID  string         `json:"entity_id"`
	Meta      map[string]any `json:"meta"`
	At        time.Time      `json:"at"`
}

type Stream struct {
	mu     sync.RWMutex
	events []Event
}

func NewStream() *Stream {
	return &Stream{events: []Event{}}
}

func (s *Stream) Append(e Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	s.events = append(s.events, e)
}

func (s *Stream) List() []Event {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Event, len(s.events))
	copy(out, s.events)
	return out
}
