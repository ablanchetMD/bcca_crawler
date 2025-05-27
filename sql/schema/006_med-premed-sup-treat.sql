-- +goose Up

CREATE TABLE medications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL UNIQUE,
  description TEXT NOT NULL DEFAULT '',
  alternate_names TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
  category TEXT NOT NULL DEFAULT ''
);

CREATE TABLE medication_prescription (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  medication UUID NOT NULL REFERENCES medications(id) ON DELETE CASCADE,
  dose TEXT NOT NULL DEFAULT '',
  route prescription_route_enum NOT NULL DEFAULT 'unknown',
  frequency TEXT NOT NULL DEFAULT '',
  duration TEXT NOT NULL DEFAULT '',
  instructions TEXT NOT NULL DEFAULT '',
  renewals INT NOT NULL DEFAULT 0,
  UNIQUE (medication, dose, route, frequency, duration, instructions)
);

CREATE TABLE protocol_treatment (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
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
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  category TEXT NOT NULL, -- Hepatic Impairment, Renal Impairment, etc.
  subcategory TEXT NOT NULL,
  adjustment TEXT NOT NULL,
  medication_id UUID NOT NULL REFERENCES medications(id) ON DELETE CASCADE,
  UNIQUE (category, subcategory, adjustment, medication_id)
);

-- +goose Down
DROP TABLE medication_modifications;
DROP TABLE medication_prescription;
DROP TABLE protocol_treatment;
DROP TABLE medications;