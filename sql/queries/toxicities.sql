-- name: AddToxicity :one
INSERT INTO toxicities (title, category, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetToxicityByID :one
SELECT
    t.id AS id,
    t.created_at AS created_at,
    t.updated_at AS updated_at,
    t.title AS title,
    t.category AS category,
    t.description AS description,
    COALESCE(
        json_agg(
            json_build_object(
                'id', tg.id,
                'created_at',
                'updated_at', 
                'grade', tg.grade,
                'description', tg.description
            )
        ) FILTER (WHERE tg.id IS NOT NULL),
        '[]'
    )::jsonb AS grades
FROM
    toxicities t
LEFT JOIN
    toxicity_grades tg ON t.id = tg.toxicity_id
WHERE t.id = $1
GROUP BY
    t.id, t.created_at, t.updated_at, t.title, t.category, t.description;

-- name: GetToxicitiesWithGrades :many
SELECT
    t.id AS id,
    t.created_at AS created_at,
    t.updated_at AS updated_at,
    t.title AS title,
    t.category AS category,
    t.description AS description,    
    COALESCE(
        json_agg(
            json_build_object(
                'id', tg.id,
                'created_at',
                'updated_at',
                'grade', tg.grade,
                'description', tg.description
            )
        ) FILTER (WHERE tg.id IS NOT NULL),
        '[]'
    )::jsonb AS grades
FROM
    toxicities t
LEFT JOIN
    toxicity_grades tg ON t.id = tg.toxicity_id
GROUP BY
    t.id, t.created_at, t.updated_at, t.title, t.category, t.description;

-- name: GetToxicitiesWithGradesAndAdjustments :many
SELECT 
    t.id,
    t.created_at,
    t.updated_at,
    t.title,
    t.category,
    t.description,
    COALESCE(
    JSON_AGG(
        JSON_BUILD_OBJECT(
            'id', tg.id,
            'created_at', tg.created_at,
            'updated_at', tg.updated_at,
            'grade', tg.grade,
            'description', tg.description,
            'adjustment', ptm.adjustment
        ) 
        ORDER BY tg.grade
        ) FILTER (WHERE tg.id IS NOT NULL),
        '[]'
    )::jsonb as grades
FROM toxicities t
LEFT JOIN toxicity_grades tg ON t.id = tg.toxicity_id
LEFT JOIN protocol_tox_modifications ptm ON tg.id = ptm.toxicity_grade_id AND ptm.protocol_id = $1
GROUP BY t.id, t.created_at, t.updated_at, t.title, t.category, t.description
ORDER BY t.title;

-- name: UpdateToxicity :one
UPDATE toxicities
SET
    updated_at = NOW(),
    title = $2,
    category = $3,
    description = $4
WHERE id = $1
RETURNING *;

-- name: UpsertToxicity :one
INSERT INTO toxicities (id,title,category,description)
VALUES (
    $1,
    $2,
    $3,
    $4    
)
ON CONFLICT (id) DO UPDATE SET
    title = EXCLUDED.title,
    category = EXCLUDED.category,
    description = EXCLUDED.description,
    updated_at = NOW()
RETURNING *;

-- name: UpsertToxicityWithGrades :exec
WITH upsert_toxicity AS (
  INSERT INTO toxicities (id, title, category, description)
  VALUES ($1, $2, $3, $4)
  ON CONFLICT (id) DO UPDATE
  SET title = EXCLUDED.title,
      category = EXCLUDED.category,
      description = EXCLUDED.description,
      updated_at = NOW()
  RETURNING id
), upsert_grades AS (
  INSERT INTO toxicity_grades (id, grade, description, toxicity_id)
  SELECT 
    unnest($5::uuid[]),
    unnest($6::grade_enum[]),
    unnest($7::text[]),
    (SELECT id FROM upsert_toxicity)
  ON CONFLICT (id) DO UPDATE
  SET grade = EXCLUDED.grade,
      description = EXCLUDED.description,
      updated_at = NOW()
)
SELECT 1;

-- name: RemoveToxicity :exec
DELETE FROM toxicities
WHERE id = $1;

-- name: AddToxicityGrade :one
INSERT INTO toxicity_grades (grade, description, toxicity_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetToxicityGrade :one
SELECT * FROM toxicity_grades
WHERE id = $1;

-- name: GetToxicityGradeByGrade :one
SELECT * FROM toxicity_grades
WHERE grade = $1 AND toxicity_id = $2;

-- name: UpdateToxicityGrade :one
UPDATE toxicity_grades
SET
    updated_at = NOW(),
    grade = $2,
    description = $3,
    toxicity_id = $4
WHERE id = $1
RETURNING *;

-- name: RemoveToxicityGrade :exec
DELETE FROM toxicity_grades
WHERE id = $1;

-- name: AddToxicityModification :one
INSERT INTO protocol_tox_modifications (adjustment, toxicity_grade_id, protocol_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpsertToxicityToProtocol :one
INSERT INTO protocol_tox_modifications (id, toxicity_grade_id, adjustment, protocol_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
SET
    toxicity_grade_id = EXCLUDED.toxicity_grade_id,
    adjustment = EXCLUDED.adjustment,
    protocol_id = EXCLUDED.protocol_id,
    updated_at = NOW()
RETURNING *;

-- name: GetToxicityModificationByProtocol :many
SELECT 
    protocol_tox_modifications.*, 
    toxicities.title AS toxicity_title, 
    toxicity_grades.description AS toxicity_grade_description, 
    toxicity_grades.grade AS toxicity_grade
FROM protocol_tox_modifications
JOIN toxicity_grades ON protocol_tox_modifications.toxicity_grade_id = toxicity_grades.id
JOIN toxicities ON toxicity_grades.toxicity_id = toxicities.id
WHERE protocol_tox_modifications.protocol_id = $1;

-- name: UpdateToxicityModification :one
UPDATE protocol_tox_modifications
SET
    updated_at = NOW(),
    adjustment = $2,
    toxicity_grade_id = $3,
    protocol_id = $4
WHERE id = $1
RETURNING *;

-- name: RemoveToxicityModification :exec
DELETE FROM protocol_tox_modifications
WHERE id = $1;

-- name: GetToxicityModificationsByTreatment :many
SELECT * FROM protocol_tox_modifications
WHERE protocol_id = $1;

-- name: GetToxicityByName :one
SELECT * FROM toxicities
WHERE title = $1;

