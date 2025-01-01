-- name: AddToxicity :one
INSERT INTO toxicities (title, category, description)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetToxicity :one
SELECT * FROM toxicities
WHERE id = $1;

-- name: UpdateToxicity :one
UPDATE toxicities
SET
    updated_at = NOW(),
    title = $2,
    category = $3,
    description = $4
WHERE id = $1
RETURNING *;

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

