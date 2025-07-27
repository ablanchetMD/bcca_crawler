-- +goose Up
CREATE TABLE protocol_tests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,  
  category TEXT NOT NULL DEFAULT 'unknown',
  comments TEXT NOT NULL DEFAULT '',
  position INT NOT NULL DEFAULT 0,
  UNIQUE (protocol_id, category)
 
);

CREATE TABLE protocol_meds (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),
  category TEXT NOT NULL DEFAULT 'unknown',
  comments TEXT NOT NULL DEFAULT '',
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  UNIQUE (category,protocol_id)
);

CREATE TABLE protocol_meds_values (
  protocol_meds_id UUID NOT NULL REFERENCES protocol_meds(id) ON DELETE CASCADE,
  medication_prescription_id UUID NOT NULL REFERENCES medication_prescription(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_meds_id, medication_prescription_id)  
);

CREATE TABLE protocol_tests_value (
  protocol_tests_id UUID NOT NULL REFERENCES protocol_tests(id) ON DELETE CASCADE,
  tests_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_tests_id, tests_id)  
);


-- +goose Down


DROP TABLE protocol_meds_values;
DROP TABLE protocol_tests_value;
DROP TABLE protocol_tests;
DROP TABLE protocol_meds;