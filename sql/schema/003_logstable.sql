-- +goose Up

CREATE Table logs (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  ip_address TEXT NOT NULL,
  user_agent TEXT NOT NULL,
  action TEXT NOT NULL
);

-- +goose Down
DROP TABLE logs;
