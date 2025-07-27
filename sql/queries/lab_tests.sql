-- name: AddTest :one
INSERT INTO tests (name, description, form_url, unit, lower_limit, upper_limit, test_category)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateTest :one
UPDATE tests
SET name = $2, description = $3, form_url = $4, unit = $5, lower_limit = $6, upper_limit = $7, test_category = $8
WHERE id = $1
RETURNING *;

-- name: UpsertTest :one
WITH input_values(id, name, description, form_url, unit, lower_limit, upper_limit, test_category) AS (
  VALUES (
    CASE 
      WHEN @id = '00000000-0000-0000-0000-000000000000'::uuid 
      THEN gen_random_uuid() 
      ELSE @id 
    END,
    @name,
    @description,
    @form_url,
    @unit,
    @lower_limit::FLOAT,
    @upper_limit::FLOAT,
    @test_category
  )
)
INSERT INTO tests (id, name, description, form_url, unit, lower_limit, upper_limit, test_category)
SELECT id, name, description, form_url, unit, lower_limit, upper_limit, test_category FROM input_values
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    form_url = EXCLUDED.form_url,
    unit = EXCLUDED.unit,
    lower_limit = EXCLUDED.lower_limit,
    upper_limit = EXCLUDED.upper_limit,
    test_category = EXCLUDED.test_category,
    updated_at = NOW()
RETURNING *;

-- name: UpsertProtoTestCategory :one
WITH input_values(id, protocol_id, category, comments,position) AS (
  VALUES (
    CASE 
      WHEN @id::uuid  = '00000000-0000-0000-0000-000000000000'
      THEN gen_random_uuid() 
      ELSE @id::uuid 
    END,
    @protocol_id::uuid,
    @category::TEXT,
    @comments::text,
    COALESCE(@position::INT, 0)
  )
)
INSERT INTO protocol_tests (id, protocol_id, category, comments,position)
SELECT id, protocol_id, category, comments,position FROM input_values
ON CONFLICT (id) DO UPDATE
SET protocol_id = EXCLUDED.protocol_id,
    category = EXCLUDED.category,
    comments = EXCLUDED.comments,
    position = EXCLUDED.position,   
    updated_at = NOW()
RETURNING *;

-- name: AddTestToProtoTestCategory :exec
INSERT INTO protocol_tests_value (protocol_tests_id, tests_id)
VALUES ($1, $2);

-- name: RemoveTestToProtoTestCategory :exec
DELETE FROM protocol_tests_value
WHERE protocol_tests_id = $1 AND tests_id = $2;

-- name: GetProtocolTests :one
SELECT COALESCE(jsonb_agg(protocol_tests_data ORDER BY protocol_tests_data.position), '[]'::jsonb) AS data
FROM (
  SELECT 
    pt.id,
    pt.created_at,
    pt.updated_at,
    pt.category,
    pt.comments,
    pt.position,
    COALESCE(tests.tests, '[]'::jsonb) AS tests
  FROM protocol_tests pt
  LEFT JOIN LATERAL (
    SELECT jsonb_agg(
      jsonb_build_object(
        'id', t.id,
        'name', t.name,
        'created_at', t.created_at,
        'updated_at', t.updated_at,
        'description', t.description,
        'form_url', t.form_url,
        'unit', t.unit,
        'lower_limit', t.lower_limit,
        'upper_limit', t.upper_limit,
        'test_category', t.test_category
      ) ORDER BY t.test_category, t.name
    ) AS tests
    FROM tests t
    JOIN protocol_tests_value ptv
      ON ptv.tests_id = t.id
    WHERE ptv.protocol_tests_id = pt.id
  ) tests ON TRUE
  WHERE pt.protocol_id = $1
) protocol_tests_data;

-- name: DeleteTest :exec
DELETE FROM tests WHERE id = $1;

-- name: GetTests :many
SELECT * FROM tests;

-- name: GetTestByID :one
SELECT * FROM tests WHERE id = $1;

-- name: GetTestsByCategory :many
SELECT * FROM tests WHERE test_category = $1;

-- name: RemoveTestCategoryByID :exec
DELETE FROM protocol_tests
WHERE id = $1;

-- name: GetTestCategoryByID :one
SELECT jsonb_build_object(
  'id', pt.id,
  'created_at', pt.created_at,
  'updated_at', pt.updated_at,
  'category', pt.category,
  'comments', pt.comments,
  'position', pt.position,
  'tests', COALESCE((
    SELECT jsonb_agg(jsonb_build_object(
      'id', t.id,
      'name', t.name,
      'created_at', t.created_at,
      'updated_at', t.updated_at,
      'description', t.description,
      'form_url', t.form_url,
      'unit', t.unit,
      'lower_limit', t.lower_limit,
      'upper_limit', t.upper_limit,
      'test_category', t.test_category
    ) ORDER BY t.test_category, t.name)
    FROM tests t
    JOIN protocol_tests_value tc ON tc.tests_id = t.id AND tc.protocol_tests_id = pt.id
  ), '[]'::jsonb)
) AS data
FROM protocol_tests pt
WHERE pt.id = $1;

-- name: GetTestByName :one
SELECT * FROM tests WHERE name = $1;

