package license

import (
	"context"
	"testing"
)

func TestMemoryRepositoryCreateGet(t *testing.T) {
	t.Parallel()
	s := NewMemoryRepository()
	created, err := s.Create(context.Background(), License{ID: "1", CreatorID: "c1", Title: "t1"})
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
}

func TestNewRepositoryDefaultsToMemory(t *testing.T) {
	t.Parallel()
	repo, err := NewRepository("", "")
	if err != nil {
		t.Fatalf("expected nil err, got %v", err)
	}
	if _, ok := repo.(*MemoryRepository); !ok {
		t.Fatalf("expected memory repo, got %T", repo)
	}
}

func TestNewRepositoryPostgresRequiresDSN(t *testing.T) {
	t.Parallel()
	_, err := NewRepository("postgres", "")
	if err == nil {
		t.Fatal("expected error for missing DATABASE_URL")
	}
}
