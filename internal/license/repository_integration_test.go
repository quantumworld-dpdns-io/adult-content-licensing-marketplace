package license

import (
	"context"
	"os"
	"testing"
	"time"
)

func TestSQLRepositoryCreateListGet(t *testing.T) {
	if os.Getenv("INTEGRATION_DB") != "1" {
		t.Skip("set INTEGRATION_DB=1 to run postgres integration tests")
	}
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Skip("DATABASE_URL is required for integration test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	db, err := OpenPostgresAndEnsureSchema(ctx, dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	defer db.Close()

	repo := NewSQLRepository(db)
	in := License{
		ID:                   "it_lic_1",
		CreatorID:            "it_creator_1",
		Title:                "Integration Test License",
		Description:          "desc",
		AITraingProhibited:   true,
		BasePriceCents:       1234,
		Currency:             "USDC",
		TerritoriesISO2Codes: []string{"us", "jp", "US"},
	}

	created, err := repo.Create(ctx, in)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if created.ID != in.ID {
		t.Fatalf("unexpected created id: %s", created.ID)
	}

	got, err := repo.Get(ctx, in.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Title != in.Title {
		t.Fatalf("unexpected title: %s", got.Title)
	}
	if len(got.TerritoriesISO2Codes) == 0 {
		t.Fatalf("expected territories")
	}

	items, err := repo.List(ctx)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(items) == 0 {
		t.Fatalf("expected non-empty list")
	}
}
