package license

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

var ErrNotFound = errors.New("license not found")

type Repository interface {
	Create(ctx context.Context, l License) (License, error)
	List(ctx context.Context) ([]License, error)
	Get(ctx context.Context, id string) (License, error)
}

type Backend string

const (
	BackendMemory   Backend = "memory"
	BackendPostgres Backend = "postgres"
)

func ParseBackend(v string) Backend {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "postgres", "pg", "postgresql":
		return BackendPostgres
	default:
		return BackendMemory
	}
}

type SQLRepository struct {
	db *sql.DB
}

func NewSQLRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, l License) (License, error) {
	territories, err := json.Marshal(l.TerritoriesISO2Codes)
	if err != nil {
		return License{}, fmt.Errorf("marshal territories: %w", err)
	}
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO licenses (
			id, creator_id, title, description, ai_training_prohibited, base_price_cents, currency, territories_json
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
	`, l.ID, l.CreatorID, l.Title, l.Description, l.AITraingProhibited, l.BasePriceCents, l.Currency, string(territories))
	if err != nil {
		return License{}, err
	}
	return l, nil
}

func (r *SQLRepository) List(ctx context.Context) ([]License, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, creator_id, title, description, ai_training_prohibited, base_price_cents, currency, COALESCE(territories_json, '[]')
		FROM licenses
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []License{}
	for rows.Next() {
		var l License
		var territoriesRaw string
		if err := rows.Scan(&l.ID, &l.CreatorID, &l.Title, &l.Description, &l.AITraingProhibited, &l.BasePriceCents, &l.Currency, &territoriesRaw); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(territoriesRaw), &l.TerritoriesISO2Codes); err != nil {
			l.TerritoriesISO2Codes = []string{}
		}
		out = append(out, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *SQLRepository) Get(ctx context.Context, id string) (License, error) {
	var l License
	var territoriesRaw string
	err := r.db.QueryRowContext(ctx, `
		SELECT id, creator_id, title, description, ai_training_prohibited, base_price_cents, currency, COALESCE(territories_json, '[]')
		FROM licenses WHERE id = $1
	`, id).Scan(&l.ID, &l.CreatorID, &l.Title, &l.Description, &l.AITraingProhibited, &l.BasePriceCents, &l.Currency, &territoriesRaw)
	if errors.Is(err, sql.ErrNoRows) {
		return License{}, ErrNotFound
	}
	if err != nil {
		return License{}, err
	}
	if err := json.Unmarshal([]byte(territoriesRaw), &l.TerritoriesISO2Codes); err != nil {
		l.TerritoriesISO2Codes = []string{}
	}
	return l, nil
}
