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
      WHEN $1 = '00000000-0000-0000-0000-000000000000'::uuid 
      THEN gen_random_uuid() 
      ELSE $1 
    END,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8::test_category_enum
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

-- name: DeleteTest :exec
DELETE FROM tests WHERE id = $1;

-- name: GetTests :many
SELECT * FROM tests;

-- name: GetTestByID :one
SELECT * FROM tests WHERE id = $1;

-- name: GetTestsByCategory :many
SELECT * FROM tests WHERE test_category = $1;

-- name: GetTestByName :one
SELECT * FROM tests WHERE name = $1;

-- name: GetTestsByProtocolByCategoryAndUrgency :many
SELECT t.*
FROM tests t
JOIN protocol_tests pt ON t.id = pt.test_id
WHERE pt.protocol_id = $1 AND pt.category = $2 AND pt.urgency = $3;

-- name: AddTestToProtocolByCategoryAndUrgency :one
INSERT INTO protocol_tests (protocol_id, test_id, category, urgency)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: RemoveTestFromProtocolByCategoryAndUrgency :exec
DELETE FROM protocol_tests WHERE protocol_id = $1 AND test_id = $2 AND category = $3 AND urgency = $4;

