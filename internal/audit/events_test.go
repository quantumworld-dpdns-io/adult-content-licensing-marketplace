package audit

import (
	"context"
	"testing"
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
	if all[0].Type != "license.created" {
		t.Fatalf("unexpected event type: %s", all[0].Type)
	}
}
