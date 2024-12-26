-- +goose Up

CREATE TABLE physicians (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  site TEXT NOT NULL DEFAULT ''
);

CREATE TABLE protocol_eligibility_criteria (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  type TEXT NOT NULL, -- Inclusion, Exclusion, or Notes
  description TEXT NOT NULL
);

CREATE TABLE protocol_eligibility_criteria_values (
  PRIMARY KEY (protocol_id, criteria_id),
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  criteria_id UUID NOT NULL REFERENCES protocol_eligibility_criteria(id) ON DELETE CASCADE
);

CREATE TABLE protocol_cautions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  description TEXT NOT NULL UNIQUE
);

CREATE TABLE protocol_cautions_values (
  PRIMARY KEY (protocol_id, caution_id),
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  caution_id UUID NOT NULL REFERENCES protocol_cautions(id) ON DELETE CASCADE
);

CREATE TABLE protocol_precautions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title TEXT NOT NULL,
  description TEXT NOT NULL
);

CREATE TABLE protocol_precautions_values (
  PRIMARY KEY (protocol_id, precaution_id),
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  precaution_id UUID NOT NULL REFERENCES protocol_precautions(id) ON DELETE CASCADE
);

CREATE TABLE tests (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL UNIQUE,
  description TEXT
);

CREATE TABLE medications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name TEXT NOT NULL,
  description TEXT,
  category TEXT NOT NULL DEFAULT ''
);

CREATE TABLE protocol_ppos (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title TEXT NOT NULL,
  url TEXT NOT NULL DEFAULT '',
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE
);

CREATE TABLE protocol_premedication (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  medication UUID NOT NULL REFERENCES medications(id) ON DELETE RESTRICT,
  dose TEXT NOT NULL,
  route TEXT NOT NULL,
  frequency TEXT NOT NULL,
  duration TEXT NOT NULL,
  notes TEXT NOT NULL DEFAULT ''
);

CREATE TABLE protocol_supportive_medication (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  medication UUID NOT NULL REFERENCES medications(id) ON DELETE RESTRICT,
  dose TEXT NOT NULL,
  route TEXT NOT NULL,
  frequency TEXT NOT NULL,
  duration TEXT NOT NULL,
  notes TEXT NOT NULL DEFAULT ''
);



CREATE TABLE toxicity_modifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  title TEXT NOT NULL,
  grade TEXT NOT NULL,
  adjustement TEXT NOT NULL,
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  UNIQUE (title,grade,adjustement,protocol_id)
);

CREATE TABLE protocol_treatment (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  medication UUID NOT NULL REFERENCES medications(id) ON DELETE RESTRICT,  
  dose TEXT NOT NULL,
  route TEXT NOT NULL,
  frequency TEXT NOT NULL,
  duration TEXT NOT NULL,
  administration_guide TEXT NOT NULL DEFAULT ''  
);

CREATE TABLE treatment_modifications (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  category TEXT NOT NULL,
  description TEXT NOT NULL,
  adjustement TEXT NOT NULL,
  treatment_id UUID NOT NULL REFERENCES protocol_treatment(id) ON DELETE CASCADE,
  UNIQUE (category,description,adjustement,treatment_id)
);

CREATE TABLE protocol_cycles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    cycle TEXT NOT NULL DEFAULT '',
    cycle_duration TEXT NOT NULL DEFAULT '',
    UNIQUE (cycle,protocol_id),
    protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE
);

CREATE TABLE treatment_cycles_junction (
  protocol_treatment_id UUID NOT NULL REFERENCES protocol_treatment(id) ON DELETE CASCADE,
  protocol_cycles_id UUID NOT NULL REFERENCES protocol_cycles(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_treatment_id, protocol_cycles_id)
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
  joi TEXT NOT NULL DEFAULT ''
);

CREATE TABLE protocol_references_value (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  reference_id UUID NOT NULL REFERENCES article_references(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, reference_id)
);


-- +goose Down

DROP TABLE physicians;
DROP TABLE protocol_eligibility_criteria_values;
DROP TABLE protocol_cautions_values;
DROP TABLE protocol_precautions_values;
DROP TABLE treatment_cycles_junction;

DROP TABLE protocol_cautions;
DROP TABLE protocol_precautions;
DROP TABLE tests;
DROP TABLE protocol_ppos;
DROP TABLE protocol_premedication;
DROP TABLE protocol_supportive_medication;
DROP TABLE treatment_modifications;
DROP TABLE toxicity_modifications;
DROP TABLE protocol_treatment;
DROP TABLE protocol_cycles;
DROP TABLE protocol_references_value;
DROP TABLE article_references;
DROP TABLE protocol_eligibility_criteria;
DROP TABLE medications;
