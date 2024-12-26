-- +goose Up
ALTER TABLE protocol_eligibility_criteria
ADD Constraint unique_description UNIQUE (description);

ALTER TABLE protocol_precautions
ADD CONSTRAINT unique_title_description UNIQUE (title, description);

ALTER TABLE medications
ADD CONSTRAINT unique_name UNIQUE (name);

ALTER TABLE protocol_ppos
ADD CONSTRAINT unique_url UNIQUE (url);

CREATE TABLE protocol_pre_medications_values (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  pre_medication_id UUID NOT NULL REFERENCES protocol_premedication(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, pre_medication_id)
);

ALTER TABLE protocol_treatment
ADD CONSTRAINT unique_combo UNIQUE (medication, dose, route, frequency, duration);

-- +goose Down

ALTER TABLE protocol_eligibility_criteria
DROP Constraint unique_description;

ALTER TABLE protocol_precautions
DROP CONSTRAINT unique_title_description;

ALTER TABLE medications
DROP CONSTRAINT unique_name;

ALTER TABLE protocol_ppos
DROP CONSTRAINT unique_url;

DROP TABLE protocol_pre_medications_values;

ALTER TABLE protocol_treatment
DROP CONSTRAINT unique_combo;
