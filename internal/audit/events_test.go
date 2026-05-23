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
}
