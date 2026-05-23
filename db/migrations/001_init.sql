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
  id UUID PRIMARY KEY,
  creator_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  ai_training_prohibited BOOLEAN NOT NULL DEFAULT true,
  base_price_cents BIGINT NOT NULL CHECK (base_price_cents >= 0),
  currency TEXT NOT NULL DEFAULT 'USD',
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS license_territories (
  license_id UUID NOT NULL REFERENCES licenses(id) ON DELETE CASCADE,
  iso_code TEXT NOT NULL,
  PRIMARY KEY (license_id, iso_code)
);
