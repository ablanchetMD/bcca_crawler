-- +goose Up

CREATE TABLE medications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL,
  category TEXT NOT NULL DEFAULT ''
);

CREATE TABLE medication_prescription (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  medication UUID NOT NULL REFERENCES medications(id) ON DELETE CASCADE,
  dose TEXT NOT NULL,
  route TEXT NOT NULL,
  frequency TEXT NOT NULL,
  duration TEXT NOT NULL,
  instructions TEXT NOT NULL DEFAULT '',
  renewals INT NOT NULL DEFAULT 0,
  UNIQUE (medication, dose, route, frequency, duration, instructions)
);

CREATE TABLE protocol_treatment (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  medication UUID NOT NULL REFERENCES medications(id) ON DELETE CASCADE,  
  dose TEXT NOT NULL,
  route TEXT NOT NULL,
  frequency TEXT NOT NULL,
  duration TEXT NOT NULL,
  administration_guide TEXT NOT NULL DEFAULT '',  
  UNIQUE (medication, dose, route, frequency, duration)
);

CREATE TABLE medication_modifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  category TEXT NOT NULL, -- Hepatic Impairment, Renal Impairment, etc.
  description TEXT NOT NULL,
  adjustment TEXT NOT NULL,
  medication_id UUID NOT NULL REFERENCES medications(id) ON DELETE CASCADE,
  UNIQUE (category, description, adjustment, medication_id)
);

-- +goose Down
DROP TABLE medication_modifications;
DROP TABLE medication_prescription;
DROP TABLE protocol_treatment;
DROP TABLE medications;