package audit

import (
	"context"
	"testing"
	"time"
)

func TestStoreAppendList(t *testing.T) {
	t.Parallel()
	s := NewMemoryStore()
	if err := s.Append(context.Background(), Event{Type: "license.created", EntityID: "lic_1"}); err != nil {
		t.Fatalf("append: %v", err)
	}
	all, err := s.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(all) != 1 {
		t.Fatalf("expected 1 event, got %d", len(all))
	}
}

func TestStoreQueryFilterAndPagination(t *testing.T) {
	t.Parallel()
	s := NewMemoryStore()
	now := time.Now().UTC()
	_ = s.Append(context.Background(), Event{Type: "license.created", EntityID: "A", At: now.Add(-2 * time.Hour)})
	_ = s.Append(context.Background(), Event{Type: "license.updated", EntityID: "A", At: now.Add(-1 * time.Hour)})
	_ = s.Append(context.Background(), Event{Type: "license.created", EntityID: "B", At: now})

	res, err := s.Query(context.Background(), Query{Type: "license.created", Limit: 10})
	if err != nil {
		t.Fatalf("query: %v", err)
	}
	if len(res) != 2 {
		t.Fatalf("expected 2 created events, got %d", len(res))
	}

	since := now.Add(-90 * time.Minute)
	res2, err := s.Query(context.Background(), Query{EntityID: "A", Since: &since, Limit: 10})
	if err != nil {
		t.Fatalf("query2: %v", err)
	}
	if len(res2) != 1 {
		t.Fatalf("expected 1 recent event for A, got %d", len(res2))
	}

	res3, err := s.Query(context.Background(), Query{Limit: 1, Offset: 1})
	if err != nil {
		t.Fatalf("query3: %v", err)
	}
	if len(res3) != 1 {
		t.Fatalf("expected 1 paged event, got %d", len(res3))
	}

	cnt, err := s.Count(context.Background(), Query{Type: "license.created"})
	if err != nil {
		t.Fatalf("count: %v", err)
	}
	if cnt != 2 {
		t.Fatalf("expected count 2, got %d", cnt)
	}
}

func TestStoreCursorPagination(t *testing.T) {
	t.Parallel()
	s := NewMemoryStore()
	now := time.Now().UTC()
	for i := 0; i < 4; i++ {
		_ = s.Append(context.Background(), Event{Type: "license.created", EntityID: "C", At: now.Add(time.Duration(i) * time.Second)})
	}
	q1 := Query{Type: "license.created", Limit: 2}
	page1, err := s.Query(context.Background(), q1)
	if err != nil {
		t.Fatalf("page1: %v", err)
	}
	if len(page1) != 2 {
		t.Fatalf("expected 2 page1 events, got %d", len(page1))
	}
	cur, err := s.NextCursor(context.Background(), q1)
	if err != nil {
		t.Fatalf("next cursor: %v", err)
	}
	if cur == "" {
		t.Fatal("expected next cursor")
	}
	page2, err := s.Query(context.Background(), Query{Type: "license.created", Limit: 2, Cursor: cur})
	if err != nil {
		t.Fatalf("page2: %v", err)
	}
	if len(page2) != 2 {
		t.Fatalf("expected 2 page2 events, got %d", len(page2))
	}
}
