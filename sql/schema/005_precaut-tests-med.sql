-- +goose Up

CREATE TABLE protocol_precautions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title TEXT NOT NULL,
  description TEXT NOT NULL,
  UNIQUE (title, description)
);

CREATE TABLE tests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL UNIQUE,
  description TEXT
);

CREATE TABLE protocol_ppos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title TEXT NOT NULL,
  url TEXT NOT NULL DEFAULT '',
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE protocol_precautions;
DROP TABLE tests;
DROP TABLE protocol_ppos;