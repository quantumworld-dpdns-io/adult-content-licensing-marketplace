package config

import (
	"os"
)

const (
	defaultPort                = "8080"
	defaultLicenseStoreBackend = "memory"
)

type Config struct {
	Port                string
	LicenseStoreBackend string
	DatabaseURL         string
	AuthSecret          string
}

func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	backend := os.Getenv("LICENSE_STORE_BACKEND")
	if backend == "" {
		backend = defaultLicenseStoreBackend
	}
	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}

	return Config{
		Port:                port,
		LicenseStoreBackend: backend,
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		AuthSecret:          secret,
	}
}
