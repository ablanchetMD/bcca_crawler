-- +goose Up


CREATE TYPE grade_enum AS ENUM ('1', '2', '3', '4','unknown');
CREATE TYPE eligibility_enum AS ENUM ('inclusion','exclusion','notes','unknown');
CREATE TYPE prescription_route_enum AS ENUM ('oral','iv','im','sc','topical','inhalation','unknown');
CREATE TYPE physician_site_enum AS ENUM ('vancouver','victoria','kelowna','surrey','prince_george','abbotsford','nanaimo','unknown');
CREATE TYPE protocol_info AS (protocol_id uuid, code text);
CREATE TYPE tumor_group_enum AS ENUM ('breast', 'lung', 'gastrointestinal', 'genitourinary', 'head_and_neck', 'gynecology', 'sarcoma', 'leukemia','bmt', 'neuro-oncology','ocular','skin','unknown_primary','lymphoma','myeloma', 'unknown');
CREATE TYPE med_adj_category_enum AS ENUM ('hepatic_impairment', 'renal_impairment','unknown');

-- +goose Down

DROP TYPE IF EXISTS eligibility_enum CASCADE;
DROP TYPE IF EXISTS physician_site_enum CASCADE;
DROP TYPE IF EXISTS prescription_route_enum CASCADE;
DROP TYPE IF EXISTS grade_enum CASCADE;
DROP TYPE IF EXISTS protocol_info CASCADE;
DROP TYPE IF EXISTS med_adj_category_enum CASCADE;
DROP TYPE IF EXISTS tumor_group_enum CASCADE;