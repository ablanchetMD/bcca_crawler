-- +goose Up

ALTER TABLE cancers 
ADD COLUMN tumor_group 
TEXT NOT NULL DEFAULT '';

-- +goose Down

ALTER TABLE cancers
DROP COLUMN tumor_group;
