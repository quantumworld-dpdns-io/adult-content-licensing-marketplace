package audit

import (
	"context"
	"os"
	"testing"
)

func TestSQLStoreAppendList(t *testing.T) {
	if os.Getenv("INTEGRATION_DB") != "1" {
		t.Skip("set INTEGRATION_DB=1 to run postgres integration tests")
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is required for integration test")
	}

	store, err := NewStore("postgres", dsn)
	if err != nil {
		t.Fatalf("new postgres store: %v", err)
	}

	if err := store.Append(context.Background(), Event{Type: "license.created", ActorSub: "u", ActorRole: "admin", EntityID: "e1", Meta: map[string]any{"k": "v"}}); err != nil {
		t.Fatalf("append: %v", err)
	}
	items, err := store.List(context.Background())
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("expected events")
	}
}
