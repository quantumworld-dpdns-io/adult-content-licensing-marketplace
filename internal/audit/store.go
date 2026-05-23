package audit

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

type Store interface {
	Append(ctx context.Context, e Event) error
	List(ctx context.Context) ([]Event, error)
}

type MemoryStore struct {
	mu     sync.RWMutex
	events []Event
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{events: []Event{}}
}

func (s *MemoryStore) Append(_ context.Context, e Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	s.events = append(s.events, e)
	return nil
}

func (s *MemoryStore) List(_ context.Context) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Event, len(s.events))
	copy(out, s.events)
	return out, nil
}

type SQLStore struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) Append(ctx context.Context, e Event) error {
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	meta, err := json.Marshal(e.Meta)
	if err != nil {
		return fmt.Errorf("marshal meta: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `
		INSERT INTO audit_events (type, actor_sub, actor_role, entity_id, meta_json, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)
	`, e.Type, e.ActorSub, e.ActorRole, e.EntityID, string(meta), e.At)
	return err
}

func (s *SQLStore) List(ctx context.Context) ([]Event, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT type, actor_sub, actor_role, entity_id, COALESCE(meta_json, '{}'), created_at
		FROM audit_events ORDER BY id DESC LIMIT 500
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Event{}
	for rows.Next() {
		var e Event
		var metaRaw string
		if err := rows.Scan(&e.Type, &e.ActorSub, &e.ActorRole, &e.EntityID, &metaRaw, &e.At); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(metaRaw), &e.Meta)
		if e.Meta == nil {
			e.Meta = map[string]any{}
		}
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

const auditSchemaSQL = `
CREATE TABLE IF NOT EXISTS audit_events (
  id BIGSERIAL PRIMARY KEY,
  type TEXT NOT NULL,
  actor_sub TEXT NOT NULL,
  actor_role TEXT NOT NULL,
  entity_id TEXT NOT NULL,
  meta_json TEXT NOT NULL DEFAULT '{}',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
`

func NewStore(backend string, databaseURL string) (Store, error) {
	if strings.ToLower(strings.TrimSpace(backend)) != "postgres" {
		return NewMemoryStore(), nil
	}
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required for postgres audit backend")
	}
	if !driverRegistered("postgres") {
		return nil, fmt.Errorf("postgres driver not registered for audit store")
	}
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if _, err := db.ExecContext(ctx, auditSchemaSQL); err != nil {
		_ = db.Close()
		return nil, err
	}
	return NewSQLStore(db), nil
}

func driverRegistered(name string) bool {
	for _, d := range sql.Drivers() {
		if d == name {
			return true
		}
	}
	return false
}
