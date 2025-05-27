-- +goose Up

CREATE TABLE physicians (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  site physician_site_enum NOT NULL DEFAULT 'unknown'
);

CREATE TABLE protocol_eligibility_criteria (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  type eligibility_enum NOT NULL DEFAULT 'unknown',
  description TEXT NOT NULL,
  unique (type, description)
);

CREATE TABLE protocol_cautions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  description TEXT NOT NULL UNIQUE
);

-- +goose Down

DROP TABLE physicians;
DROP TABLE protocol_eligibility_criteria;
DROP TABLE protocol_cautions;
