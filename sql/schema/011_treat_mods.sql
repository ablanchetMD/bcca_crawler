-- +goose Up

CREATE TABLE toxicities (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title TEXT NOT NULL UNIQUE,
  category TEXT NOT NULL,
  description TEXT NOT NULL  
);

CREATE TABLE toxicity_grades (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  grade grade_enum NOT NULL DEFAULT 'unknown',
  description TEXT NOT NULL,
  toxicity_id UUID NOT NULL REFERENCES toxicities(id) ON DELETE CASCADE,
  UNIQUE (grade, toxicity_id)
);

CREATE TABLE protocol_tox_modifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    adjustment TEXT NOT NULL,
    toxicity_grade_id UUID NOT NULL REFERENCES toxicity_grades(id) ON DELETE CASCADE,
    protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
    UNIQUE (adjustment, toxicity_grade_id, protocol_id)
);

-- +goose Down

DROP TABLE protocol_tox_modifications;
DROP TABLE toxicity_grades;
DROP TABLE toxicities;

