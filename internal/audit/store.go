package audit

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Query struct {
	Type     string
	EntityID string
	Since    *time.Time
	Until    *time.Time
	Limit    int
	Offset   int
	Cursor   string
}

type Cursor struct {
	At time.Time `json:"at"`
	ID int64     `json:"id"`
}

func EncodeCursor(c Cursor) string {
	b, _ := json.Marshal(c)
	return base64.RawURLEncoding.EncodeToString(b)
}

func DecodeCursor(v string) (Cursor, error) {
	var c Cursor
	if strings.TrimSpace(v) == "" {
		return c, nil
	}
	raw, err := base64.RawURLEncoding.DecodeString(v)
	if err != nil {
		return c, fmt.Errorf("invalid cursor")
	}
	if err := json.Unmarshal(raw, &c); err != nil {
		return c, fmt.Errorf("invalid cursor")
	}
	if c.At.IsZero() || c.ID <= 0 {
		return c, fmt.Errorf("invalid cursor")
	}
	return c, nil
}

func (q Query) Normalized() Query {
	out := q
	if out.Limit <= 0 || out.Limit > 500 {
		out.Limit = 100
	}
	if out.Offset < 0 {
		out.Offset = 0
	}
	out.Type = strings.TrimSpace(out.Type)
	out.EntityID = strings.TrimSpace(out.EntityID)
	out.Cursor = strings.TrimSpace(out.Cursor)
	return out
}

type Store interface {
	Append(ctx context.Context, e Event) error
	List(ctx context.Context) ([]Event, error)
	Query(ctx context.Context, q Query) ([]Event, error)
	Count(ctx context.Context, q Query) (int, error)
	NextCursor(ctx context.Context, q Query) (string, error)
}

type memoryEvent struct {
	ID    int64
	Event Event
}

type MemoryStore struct {
	mu     sync.RWMutex
	events []memoryEvent
	nextID int64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{events: []memoryEvent{}, nextID: 1}
}

func (s *MemoryStore) Append(_ context.Context, e Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	s.events = append(s.events, memoryEvent{ID: s.nextID, Event: e})
	s.nextID++
	return nil
}

func (s *MemoryStore) List(ctx context.Context) ([]Event, error) {
	return s.Query(ctx, Query{})
}

func (s *MemoryStore) Query(_ context.Context, q Query) ([]Event, error) {
	q = q.Normalized()
	s.mu.RLock()
	defer s.mu.RUnlock()
	filtered, err := s.filterLocked(q)
	if err != nil {
		return nil, err
	}
	if q.Cursor == "" {
		if q.Offset >= len(filtered) {
			return []Event{}, nil
		}
		end := q.Offset + q.Limit
		if end > len(filtered) {
			end = len(filtered)
		}
		out := make([]Event, 0, end-q.Offset)
		for _, me := range filtered[q.Offset:end] {
			out = append(out, me.Event)
		}
		return out, nil
	}
	end := q.Limit
	if end > len(filtered) {
		end = len(filtered)
	}
	out := make([]Event, 0, end)
	for _, me := range filtered[:end] {
		out = append(out, me.Event)
	}
	return out, nil
}

func (s *MemoryStore) Count(_ context.Context, q Query) (int, error) {
	q = q.Normalized()
	s.mu.RLock()
	defer s.mu.RUnlock()
	filtered, err := s.filterLocked(q)
	if err != nil {
		return 0, err
	}
	return len(filtered), nil
}

func (s *MemoryStore) NextCursor(_ context.Context, q Query) (string, error) {
	q = q.Normalized()
	s.mu.RLock()
	defer s.mu.RUnlock()
	filtered, err := s.filterLocked(q)
	if err != nil {
		return "", err
	}
	if len(filtered) <= q.Limit {
		return "", nil
	}
	last := filtered[q.Limit-1]
	return EncodeCursor(Cursor{At: last.Event.At, ID: last.ID}), nil
}

func (s *MemoryStore) filterLocked(q Query) ([]memoryEvent, error) {
	var cursor Cursor
	var err error
	if q.Cursor != "" {
		cursor, err = DecodeCursor(q.Cursor)
		if err != nil {
			return nil, err
		}
	}
	filtered := make([]memoryEvent, 0, len(s.events))
	for i := len(s.events) - 1; i >= 0; i-- {
		me := s.events[i]
		e := me.Event
		if q.Type != "" && e.Type != q.Type {
			continue
		}
		if q.EntityID != "" && e.EntityID != q.EntityID {
			continue
		}
		if q.Since != nil && e.At.Before(*q.Since) {
			continue
		}
		if q.Until != nil && e.At.After(*q.Until) {
			continue
		}
		if q.Cursor != "" {
			if e.At.After(cursor.At) {
				continue
			}
			if e.At.Equal(cursor.At) && me.ID >= cursor.ID {
				continue
			}
		}
		filtered = append(filtered, me)
	}
	return filtered, nil
}

type SQLStore struct{ db *sql.DB }

func NewSQLStore(db *sql.DB) *SQLStore { return &SQLStore{db: db} }

func (s *SQLStore) Append(ctx context.Context, e Event) error {
	if e.At.IsZero() {
		e.At = time.Now().UTC()
	}
	meta, err := json.Marshal(e.Meta)
	if err != nil {
		return fmt.Errorf("marshal meta: %w", err)
	}
	_, err = s.db.ExecContext(ctx, `INSERT INTO audit_events (type, actor_sub, actor_role, entity_id, meta_json, created_at) VALUES ($1,$2,$3,$4,$5,$6)`, e.Type, e.ActorSub, e.ActorRole, e.EntityID, string(meta), e.At)
	return err
}

func (s *SQLStore) List(ctx context.Context) ([]Event, error) { return s.Query(ctx, Query{}) }

func (s *SQLStore) Query(ctx context.Context, q Query) ([]Event, error) {
	q = q.Normalized()
	where, args, err := buildWhere(q)
	if err != nil {
		return nil, err
	}
	if q.Cursor == "" {
		args = append(args, q.Limit)
		limitArg := "$" + strconv.Itoa(len(args))
		args = append(args, q.Offset)
		offsetArg := "$" + strconv.Itoa(len(args))
		stmt := `SELECT id, type, actor_sub, actor_role, entity_id, COALESCE(meta_json, '{}'), created_at FROM audit_events WHERE ` + strings.Join(where, " AND ") + ` ORDER BY created_at DESC, id DESC LIMIT ` + limitArg + ` OFFSET ` + offsetArg
		return scanEvents(s.db.QueryContext(ctx, stmt, args...))
	}
	args = append(args, q.Limit)
	limitArg := "$" + strconv.Itoa(len(args))
	stmt := `SELECT id, type, actor_sub, actor_role, entity_id, COALESCE(meta_json, '{}'), created_at FROM audit_events WHERE ` + strings.Join(where, " AND ") + ` ORDER BY created_at DESC, id DESC LIMIT ` + limitArg
	return scanEvents(s.db.QueryContext(ctx, stmt, args...))
}

func scanEvents(rows *sql.Rows, err error) ([]Event, error) {
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Event{}
	for rows.Next() {
		var id int64
		var e Event
		var metaRaw string
		if err := rows.Scan(&id, &e.Type, &e.ActorSub, &e.ActorRole, &e.EntityID, &metaRaw, &e.At); err != nil {
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

func (s *SQLStore) Count(ctx context.Context, q Query) (int, error) {
	q = q.Normalized()
	where, args, err := buildWhere(q)
	if err != nil {
		return 0, err
	}
	stmt := `SELECT COUNT(1) FROM audit_events WHERE ` + strings.Join(where, " AND ")
	var n int
	if err := s.db.QueryRowContext(ctx, stmt, args...).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

func (s *SQLStore) NextCursor(ctx context.Context, q Query) (string, error) {
	q = q.Normalized()
	where, args, err := buildWhere(q)
	if err != nil {
		return "", err
	}
	args = append(args, q.Limit)
	limitArg := "$" + strconv.Itoa(len(args))
	stmt := `SELECT id, created_at FROM audit_events WHERE ` + strings.Join(where, " AND ") + ` ORDER BY created_at DESC, id DESC LIMIT ` + limitArg
	rows, err := s.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	var ids []int64
	var ats []time.Time
	for rows.Next() {
		var id int64
		var at time.Time
		if err := rows.Scan(&id, &at); err != nil {
			return "", err
		}
		ids = append(ids, id)
		ats = append(ats, at)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	if len(ids) < q.Limit {
		return "", nil
	}
	idx := len(ids) - 1
	return EncodeCursor(Cursor{At: ats[idx], ID: ids[idx]}), nil
}

func buildWhere(q Query) ([]string, []any, error) {
	where := []string{"1=1"}
	args := []any{}
	if q.Type != "" {
		args = append(args, q.Type)
		where = append(where, "type = $"+strconv.Itoa(len(args)))
	}
	if q.EntityID != "" {
		args = append(args, q.EntityID)
		where = append(where, "entity_id = $"+strconv.Itoa(len(args)))
	}
	if q.Since != nil {
		args = append(args, *q.Since)
		where = append(where, "created_at >= $"+strconv.Itoa(len(args)))
	}
	if q.Until != nil {
		args = append(args, *q.Until)
		where = append(where, "created_at <= $"+strconv.Itoa(len(args)))
	}
	if q.Cursor != "" {
		c, err := DecodeCursor(q.Cursor)
		if err != nil {
			return nil, nil, err
		}
		args = append(args, c.At)
		aArg := "$" + strconv.Itoa(len(args))
		args = append(args, c.ID)
		iArg := "$" + strconv.Itoa(len(args))
		where = append(where, "(created_at < "+aArg+" OR (created_at = "+aArg+" AND id < "+iArg+"))")
	}
	return where, args, nil
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
CREATE INDEX IF NOT EXISTS idx_audit_events_type ON audit_events(type);
CREATE INDEX IF NOT EXISTS idx_audit_events_entity_id ON audit_events(entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_events_created_at ON audit_events(created_at);
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
