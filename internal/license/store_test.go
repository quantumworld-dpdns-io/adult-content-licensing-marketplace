package license

import "testing"

func TestStoreCreateGet(t *testing.T) {
	t.Parallel()
	s := NewStore()
	created := s.Create(License{ID: "1", CreatorID: "c1", Title: "t1"})
	if created.ID != "1" {
		t.Fatalf("expected id 1, got %s", created.ID)
	}
	got, err := s.Get("1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Title != "t1" {
		t.Fatalf("expected title t1, got %s", got.Title)
	}
}
