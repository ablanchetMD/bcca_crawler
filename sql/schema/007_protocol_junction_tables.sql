-- +goose Up

CREATE TABLE protocol_baseline_tests (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, test_id)
);

CREATE TABLE protocol_baseline_tests_non_urgent (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, test_id)
);

CREATE TABLE protocol_baseline_tests_if_necessary (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, test_id)
);

CREATE TABLE protocol_followup_tests (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, test_id)
);

CREATE TABLE protocol_followup_tests_if_necessary (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, test_id)
);

CREATE TABLE protocol_contact_physicians (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  physician_id UUID NOT NULL REFERENCES physicians(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, physician_id)
);

CREATE TABLE protocol_supportive_medications (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  supportive_medication_id UUID NOT NULL REFERENCES protocol_supportive_medication(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, supportive_medication_id)
);

CREATE TABLE protocol_cycles_values (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  cycle_id UUID NOT NULL REFERENCES protocol_cycles(id) ON DELETE CASCADE,
  PRIMARY KEY (protocol_id, cycle_id)
);

-- +goose Down

DROP TABLE protocol_baseline_tests;
DROP TABLE protocol_baseline_tests_non_urgent;
DROP TABLE protocol_baseline_tests_if_necessary;
DROP TABLE protocol_followup_tests;
DROP TABLE protocol_followup_tests_if_necessary;
DROP TABLE protocol_contact_physicians;
DROP TABLE protocol_supportive_medications;
DROP TABLE protocol_cycles_values;

