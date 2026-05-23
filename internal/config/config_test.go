package config

import "testing"

func TestLoadDefaultPort(t *testing.T) {
	t.Setenv("PORT", "")
	t.Setenv("LICENSE_STORE_BACKEND", "")
	t.Setenv("AUTH_SECRET", "")
	cfg := Load()
	if cfg.Port != defaultPort {
		t.Fatalf("expected default port %s, got %s", defaultPort, cfg.Port)
	}
	if cfg.LicenseStoreBackend != defaultLicenseStoreBackend {
		t.Fatalf("expected default backend %s, got %s", defaultLicenseStoreBackend, cfg.LicenseStoreBackend)
	}
	if cfg.AuthSecret != "dev-secret" {
		t.Fatalf("expected default auth secret")
	}
}

func TestLoadPortFromEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("LICENSE_STORE_BACKEND", "postgres")
	t.Setenv("DATABASE_URL", "postgres://u:p@localhost:5432/db")
	t.Setenv("AUTH_SECRET", "prod-secret")
	cfg := Load()
	if cfg.Port != "9090" {
		t.Fatalf("expected port 9090, got %s", cfg.Port)
	}
	if cfg.LicenseStoreBackend != "postgres" {
		t.Fatalf("expected postgres backend, got %s", cfg.LicenseStoreBackend)
	}
	if cfg.DatabaseURL == "" {
		t.Fatalf("expected database url")
	}
	if cfg.AuthSecret != "prod-secret" {
		t.Fatalf("expected auth secret from env")
	}
}
