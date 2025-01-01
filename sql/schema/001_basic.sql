-- +goose Up
CREATE TABLE protocols (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  tumor_group TEXT NOT NULL,
  code TEXT NOT NULL UNIQUE,
  name TEXT NOT NULL,
  tags TEXT[] NOT NULL DEFAULT '{}',
  notes TEXT NOT NULL DEFAULT '',
  protocol_url TEXT NOT NULL DEFAULT '',
  patient_handout_url TEXT NOT NULL DEFAULT '',
  revised_on TEXT NOT NULL DEFAULT '',
  activated_on TEXT NOT NULL DEFAULT ''
);


CREATE TABLE cancers (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  tumor_group TEXT NOT NULL DEFAULT '',  
  code TEXT,
  name TEXT,
  tags TEXT[] NOT NULL DEFAULT '{}',
  notes TEXT NOT NULL DEFAULT ''
);

CREATE TABLE cancer_protocols (
  PRIMARY KEY (cancer_id, protocol_id),
  cancer_id UUID NOT NULL REFERENCES cancers(id) ON DELETE CASCADE,
  protocol_id UUID NOT NULL REFERENCES protocols(id) ON DELETE CASCADE
);

-- +goose Down

DROP TABLE cancer_protocols;
DROP TABLE cancers;
DROP TABLE protocols;