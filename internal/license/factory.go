package license

import (
	"database/sql"
	"fmt"
)

func NewRepository(backend string, databaseURL string) (Repository, error) {
	switch ParseBackend(backend) {
	case BackendPostgres:
		if databaseURL == "" {
			return nil, fmt.Errorf("DATABASE_URL is required for postgres backend")
		}
		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			return nil, fmt.Errorf("open postgres: %w", err)
		}
		return NewSQLRepository(db), nil
	default:
		return NewMemoryRepository(), nil
	}
}
