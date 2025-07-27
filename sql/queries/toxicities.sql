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
        jsonb_agg(
            jsonb_build_object(
                'id', tg.id,
                'created_at',tg.created_at,
                'updated_at',tg.updated_at,
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
        jsonb_agg(
            jsonb_build_object(
                'id', tg.id,
                'created_at', tg.created_at,
                'updated_at', tg.updated_at,
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

-- name: GetToxicitiesWithGradesAndAdjustmentsByProtocol :many
SELECT 
    t.id,
    t.created_at,
    t.updated_at,
    t.title,
    t.category,
    t.description,
    COALESCE(
    JSONb_AGG(
        JSONb_BUILD_OBJECT(
            'id', tg.id,
            'created_at', tg.created_at,
            'updated_at', tg.updated_at,
            'grade', tg.grade,
            'description', tg.description,
            'adjustment', ptm.adjustment,
            'adjustment_id', ptm.id
        ) 
        ORDER BY tg.grade
        ) FILTER (WHERE tg.id IS NOT NULL),
        '[]'
    )::jsonb as grades
FROM toxicities t
LEFT JOIN toxicity_grades tg ON tg.toxicity_id = t.id
LEFT JOIN protocol_tox_modifications ptm 
  ON ptm.toxicity_grade_id = tg.id 
  AND ptm.protocol_id = $1
WHERE EXISTS (
  SELECT 1 FROM protocol_tox_modifications ptm2
  JOIN toxicity_grades tg2 ON ptm2.toxicity_grade_id = tg2.id
  WHERE tg2.toxicity_id = t.id AND ptm2.protocol_id = $1
)
GROUP BY t.id;


-- name: GetToxicitiesWithGradesAndAdjustments :many
SELECT 
    t.id,
    t.created_at,
    t.updated_at,
    t.title,
    t.category,
    t.description,
    COALESCE(
    JSONb_AGG(
        JSONb_BUILD_OBJECT(
            'id', tg.id,
            'created_at', tg.created_at,
            'updated_at', tg.updated_at,
            'grade', tg.grade,
            'description', tg.description,
            'adjustment', ptm.adjustment,
            'adjustment_id', ptm.id
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
    unnest(@grade_ids::uuid[]),
    unnest(@grade_number::grade_enum[]),
    unnest(@grade_description::text[]),
    (SELECT id FROM upsert_toxicity)
  ON CONFLICT (id) DO UPDATE
  SET grade = EXCLUDED.grade,
      description = EXCLUDED.description,
      updated_at = NOW()
)
SELECT 1;

-- name: UpsertToxicityModification :exec
WITH input_data AS (
    SELECT 
        unnest(@id::uuid[]) AS id,
        unnest(@grade_ids::uuid[]) AS grade_id,
        unnest(@adjustment::text[]) AS adjustment
)
INSERT INTO protocol_tox_modifications (
    id, 
    toxicity_grade_id, 
    protocol_id, 
    adjustment
)
SELECT
    id,
    grade_id,
    @protocol_id::uuid,
    adjustment
FROM input_data
ON CONFLICT (id) DO UPDATE SET
    toxicity_grade_id = EXCLUDED.toxicity_grade_id,
    adjustment = EXCLUDED.adjustment,
    updated_at = NOW();

-- name: DeleteProtocolToxModificationsByProtocolAndToxicity :exec
DELETE FROM protocol_tox_modifications ptm
USING toxicity_grades tg
WHERE ptm.toxicity_grade_id = tg.id
  AND ptm.protocol_id = @protocol_id
  AND tg.toxicity_id = @toxicity_id;

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

