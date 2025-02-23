-- +goose Up
CREATE TABLE protocol_tests (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  test_id UUID NOT NULL REFERENCES tests(id) ON DELETE CASCADE,
  category category_enum NOT NULL DEFAULT 'unknown',
  urgency urgency_enum NOT NULL DEFAULT 'unknown',  
  PRIMARY KEY (protocol_id, test_id, category, urgency)
);

CREATE TABLE protocol_meds (
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE,
  prescription_id UUID NOT NULL REFERENCES medication_prescription(id) ON DELETE CASCADE,
  category med_proto_category_enum NOT NULL DEFAULT 'unknown',  
  PRIMARY KEY (protocol_id, prescription_id, category)
);

-- +goose Down

DROP TABLE protocol_tests;
DROP TABLE protocol_meds;