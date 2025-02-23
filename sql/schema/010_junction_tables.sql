-- +goose Up

CREATE TABLE protocol_contact_physicians (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  physician_id UUID NOT NULL REFERENCES physicians(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, physician_id)
);

CREATE TABLE protocol_eligibility_criteria_values (
  PRIMARY KEY (protocol_id, criteria_id),
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  criteria_id UUID NOT NULL REFERENCES protocol_eligibility_criteria(id) ON DELETE CASCADE
);

CREATE TABLE protocol_cautions_values (
  PRIMARY KEY (protocol_id, caution_id),
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  caution_id UUID NOT NULL REFERENCES protocol_cautions(id) ON DELETE CASCADE
);

CREATE TABLE protocol_precautions_values (
  PRIMARY KEY (protocol_id, precaution_id),
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  precaution_id UUID NOT NULL REFERENCES protocol_precautions(id) ON DELETE CASCADE
);

CREATE TABLE treatment_cycles_values (
  protocol_treatment_id UUID NOT NULL REFERENCES protocol_treatment(id) ON DELETE CASCADE,
  protocol_cycles_id UUID NOT NULL REFERENCES protocol_cycles(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_treatment_id, protocol_cycles_id)
);

CREATE TABLE protocol_references_value (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  reference_id UUID NOT NULL REFERENCES article_references(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, reference_id)
);

-- +goose Down

DROP TABLE protocol_contact_physicians;
DROP TABLE protocol_eligibility_criteria_values;
DROP TABLE protocol_cautions_values;
DROP TABLE protocol_precautions_values;
DROP TABLE treatment_cycles_values;
DROP TABLE protocol_references_value;
