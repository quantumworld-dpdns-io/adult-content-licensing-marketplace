package license

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sort"
	"strings"
)

var ErrNotFound = errors.New("license not found")

type Repository interface {
	Create(ctx context.Context, l License) (License, error)
	List(ctx context.Context) ([]License, error)
	Get(ctx context.Context, id string) (License, error)
	Update(ctx context.Context, id string, in License) (License, error)
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

type SQLRepository struct{ db *sql.DB }

func NewSQLRepository(db *sql.DB) *SQLRepository { return &SQLRepository{db: db} }

func (r *SQLRepository) Create(ctx context.Context, l License) (License, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return License{}, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `INSERT INTO licenses (id, creator_id, title, description, ai_training_prohibited, base_price_cents, currency, version) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`, l.ID, l.CreatorID, l.Title, l.Description, l.AITraingProhibited, l.BasePriceCents, l.Currency, 1)
	if err != nil {
		return License{}, err
	}
	for _, iso := range dedupeTerritories(l.TerritoriesISO2Codes) {
		_, err = tx.ExecContext(ctx, `INSERT INTO license_territories (license_id, iso_code) VALUES ($1,$2) ON CONFLICT (license_id, iso_code) DO NOTHING`, l.ID, strings.ToUpper(iso))
		if err != nil {
			return License{}, err
		}
	}
	if err := tx.Commit(); err != nil {
		return License{}, err
	}
	l.Version = 1
	return l, nil
}

func (r *SQLRepository) List(ctx context.Context) ([]License, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT id, creator_id, title, description, ai_training_prohibited, base_price_cents, currency, version FROM licenses ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []License{}
	for rows.Next() {
		var l License
		if err := rows.Scan(&l.ID, &l.CreatorID, &l.Title, &l.Description, &l.AITraingProhibited, &l.BasePriceCents, &l.Currency, &l.Version); err != nil {
			return nil, err
		}
		territories, err := r.territoriesForLicense(ctx, l.ID)
		if err != nil {
			return nil, err
		}
		l.TerritoriesISO2Codes = territories
		out = append(out, l)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *SQLRepository) Get(ctx context.Context, id string) (License, error) {
	var l License
	err := r.db.QueryRowContext(ctx, `SELECT id, creator_id, title, description, ai_training_prohibited, base_price_cents, currency, version FROM licenses WHERE id=$1`, id).Scan(&l.ID, &l.CreatorID, &l.Title, &l.Description, &l.AITraingProhibited, &l.BasePriceCents, &l.Currency, &l.Version)
	if errors.Is(err, sql.ErrNoRows) {
		return License{}, ErrNotFound
	}
	if err != nil {
		return License{}, err
	}
	territories, err := r.territoriesForLicense(ctx, l.ID)
	if err != nil {
		return License{}, err
	}
	l.TerritoriesISO2Codes = territories
	return l, nil
}

func (r *SQLRepository) Update(ctx context.Context, id string, in License) (License, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return License{}, err
	}
	defer tx.Rollback()

	if in.Version <= 0 {
		return License{}, ErrBadVersion
	}
	res, err := tx.ExecContext(ctx, `UPDATE licenses SET creator_id=$2,title=$3,description=$4,ai_training_prohibited=$5,base_price_cents=$6,currency=$7,version=version+1 WHERE id=$1 AND version=$8`, id, in.CreatorID, in.Title, in.Description, in.AITraingProhibited, in.BasePriceCents, in.Currency, in.Version)
	if err != nil {
		return License{}, err
	}
	aff, _ := res.RowsAffected()
	if aff == 0 {
		var exists int
		if err := tx.QueryRowContext(ctx, `SELECT COUNT(1) FROM licenses WHERE id=$1`, id).Scan(&exists); err != nil {
			return License{}, err
		}
		if exists == 0 {
			return License{}, ErrNotFound
		}
		return License{}, ErrConflict
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM license_territories WHERE license_id=$1`, id); err != nil {
		return License{}, err
	}
	for _, iso := range dedupeTerritories(in.TerritoriesISO2Codes) {
		if _, err := tx.ExecContext(ctx, `INSERT INTO license_territories (license_id, iso_code) VALUES ($1,$2)`, id, strings.ToUpper(iso)); err != nil {
			return License{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return License{}, err
	}
	in.ID = id
	in.Version = in.Version + 1
	return in, nil
}

func (r *SQLRepository) territoriesForLicense(ctx context.Context, licenseID string) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, `SELECT iso_code FROM license_territories WHERE license_id=$1 ORDER BY iso_code`, licenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var iso string
		if err := rows.Scan(&iso); err != nil {
			return nil, err
		}
		out = append(out, iso)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func dedupeTerritories(in []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(in))
	for _, v := range in {
		vv := strings.ToUpper(strings.TrimSpace(v))
		if vv == "" {
			continue
		}
		if _, ok := seen[vv]; ok {
			continue
		}
		seen[vv] = struct{}{}
		out = append(out, vv)
	}
	sort.Strings(out)
	return out
}

func OpenPostgresAndEnsureSchema(ctx context.Context, dsn string) (*sql.DB, error) {
	if dsn == "" {
		return nil, fmt.Errorf("dsn is required")
	}
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	if err := EnsureSchema(ctx, db); err != nil {
		_ = db.Close()
		return nil, err
	}
	return db, nil
}
