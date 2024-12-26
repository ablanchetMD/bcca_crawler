-- +goose Up

ALTER TABLE protocols 
ADD COLUMN protocol_url TEXT NOT NULL DEFAULT '',
ADD COLUMN patient_handout_url TEXT NOT NULL DEFAULT '';

-- +goose Down

ALTER TABLE protocols
DROP COLUMN protocol_url,
DROP COLUMN patient_handout_url;
