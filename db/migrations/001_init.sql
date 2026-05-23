CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  email TEXT UNIQUE NOT NULL,
  role TEXT NOT NULL CHECK (role IN ('creator','buyer','admin')),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS creators (
  user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  display_name TEXT NOT NULL,
  bio TEXT DEFAULT '',
  verified BOOLEAN NOT NULL DEFAULT false
);

CREATE TABLE IF NOT EXISTS licenses (
  id TEXT PRIMARY KEY,
  creator_id TEXT NOT NULL,
  title TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  ai_training_prohibited BOOLEAN NOT NULL DEFAULT true,
  base_price_cents BIGINT NOT NULL CHECK (base_price_cents >= 0),
  currency TEXT NOT NULL DEFAULT 'USDC',
  territories_json TEXT NOT NULL DEFAULT '[]',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS license_territories (
  license_id TEXT NOT NULL REFERENCES licenses(id) ON DELETE CASCADE,
  iso_code TEXT NOT NULL,
  PRIMARY KEY (license_id, iso_code)
);
