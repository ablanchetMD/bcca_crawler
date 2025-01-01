-- +goose Up

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  email TEXT NOT NULL UNIQUE,
  role TEXT NOT NULL DEFAULT 'user',
  is_verified BOOLEAN NOT NULL DEFAULT FALSE,
  deleted_at TIMESTAMP,
  deleted_by UUID REFERENCES users(id),
  last_active TIMESTAMP,
  password TEXT NOT NULL
);

CREATE TABLE refresh_tokens (
  token TEXT PRIMARY KEY NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  expires_at TIMESTAMP NOT NULL,
  revoked_at TIMESTAMP,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE refresh_tokens;
DROP TABLE users;

