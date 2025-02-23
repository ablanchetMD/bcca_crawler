-- +goose Up

CREATE TABLE protocol_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    cycle TEXT NOT NULL DEFAULT '',
    cycle_duration TEXT NOT NULL DEFAULT '',
    UNIQUE (cycle,protocol_id),
    protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE
);

CREATE TABLE article_references (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title TEXT NOT NULL,
  authors TEXT NOT NULL,
  journal TEXT NOT NULL,
  year TEXT NOT NULL,
  UNIQUE (title,authors,journal,year),
  pmid TEXT NOT NULL DEFAULT '',
  doi TEXT NOT NULL DEFAULT ''
);

-- +goose Down

DROP TABLE protocol_cycles;
DROP TABLE article_references;
