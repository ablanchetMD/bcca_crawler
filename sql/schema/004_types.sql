-- +goose Up

CREATE TYPE category_enum AS ENUM ('baseline', 'followup','unknown');
CREATE TYPE urgency_enum AS ENUM ('urgent', 'non_urgent', 'if_necessary','unknown');
CREATE TYPE grade_enum AS ENUM ('1', '2', '3', '4','unknown');
CREATE TYPE med_proto_category_enum AS ENUM ('premed', 'support','unknown');
CREATE TYPE eligibility_enum AS ENUM ('inclusion','exclusion','notes','unknown');
CREATE TYPE prescription_route_enum AS ENUM ('oral','iv','im','sc','topical','inhalation','unknown');
CREATE TYPE physician_site_enum AS ENUM ('vancouver','victoria','kelowna','surrey','prince_george','abbotsford','nanaimo','unknown');

-- +goose Down

DROP TYPE IF EXISTS category_enum CASCADE;
DROP TYPE IF EXISTS urgency_enum CASCADE;
DROP TYPE IF EXISTS eligibility_enum CASCADE;
DROP TYPE IF EXISTS physician_site_enum CASCADE;
DROP TYPE IF EXISTS med_proto_category_enum CASCADE;
DROP TYPE IF EXISTS prescription_route_enum CASCADE;
DROP TYPE IF EXISTS grade_enum CASCADE;