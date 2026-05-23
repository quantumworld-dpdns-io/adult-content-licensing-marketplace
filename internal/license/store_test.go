package license

import (
	"context"
	"testing"
)

func TestMemoryRepositoryCreateGetUpdate(t *testing.T) {
	t.Parallel()
	s := NewMemoryRepository()
	created, err := s.Create(context.Background(), License{ID: "1", CreatorID: "c1", Title: "t1", Currency: "USDC"})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if created.ID != "1" {
		t.Fatalf("expected id 1, got %s", created.ID)
	}
	got, err := s.Get(context.Background(), "1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != "t1" {
		t.Fatalf("expected title t1, got %s", got.Title)
	}
	updated, err := s.Update(context.Background(), "1", License{CreatorID: "c1", Title: "t2", Currency: "ETH", Version: got.Version})
	if err != nil {
		t.Fatalf("unexpected update error: %v", err)
	}
	if updated.Title != "t2" || updated.ID != "1" {
		t.Fatalf("unexpected updated value: %+v", updated)
	}
}

func TestMemoryRepositoryUpdateConflict(t *testing.T) {
	t.Parallel()
	s := NewMemoryRepository()
	_, _ = s.Create(context.Background(), License{ID: "1", CreatorID: "c1", Title: "t1", Currency: "USDC"})
	_, err := s.Update(context.Background(), "1", License{CreatorID: "c1", Title: "t2", Currency: "USDC", Version: 999})
	if err != ErrConflict {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}
