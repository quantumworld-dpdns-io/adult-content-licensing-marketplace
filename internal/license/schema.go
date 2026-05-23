package license

import (
	"context"
	"database/sql"
)

const bootstrapSchemaSQL = `
CREATE TABLE IF NOT EXISTS licenses (
  id TEXT PRIMARY KEY,
  creator_id TEXT NOT NULL,
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  ai_training_prohibited BOOLEAN NOT NULL DEFAULT true,
  base_price_cents BIGINT NOT NULL CHECK (base_price_cents >= 0),
  currency TEXT NOT NULL DEFAULT 'USDC',
  version BIGINT NOT NULL DEFAULT 1,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS license_territories (
  license_id TEXT NOT NULL REFERENCES licenses(id) ON DELETE CASCADE,
  iso_code TEXT NOT NULL,
  PRIMARY KEY (license_id, iso_code)
);
`

func EnsureSchema(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, bootstrapSchemaSQL)
	return err
}
