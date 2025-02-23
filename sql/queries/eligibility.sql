-- name: InsertEligibilityCriteria :one
INSERT INTO protocol_eligibility_criteria (type, description)
VALUES ($1, $2)
RETURNING *;

-- name: UpsertEligibilityCriteria :one
INSERT INTO protocol_eligibility_criteria (id, type, description, updated_at)
VALUES ($1, $2, $3, NOW())
ON CONFLICT (id) DO UPDATE
SET type = EXCLUDED.type,
    description = EXCLUDED.description,
    updated_at = NOW()
RETURNING *;

-- name: AddEligibilityToProtocol :exec
INSERT INTO protocol_eligibility_criteria_values (protocol_id, criteria_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;


-- name: UpdateEligibilityCriteria :one
UPDATE protocol_eligibility_criteria
SET
    updated_at = NOW(),
    type = $2,
    description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteEligibilityCriteria :exec
DELETE FROM protocol_eligibility_criteria
WHERE id = $1;

-- name: GetElibilityCriteriaByDescription :one
SELECT * FROM protocol_eligibility_criteria
WHERE description = $1;

-- name: LinkEligibilityToProtocol :exec
INSERT INTO protocol_eligibility_criteria_values (protocol_id, criteria_id)
VALUES ($1, $2);

-- name: GetElibilityCriteria :many
SELECT pec.*, ARRAY_AGG(ROW(pecv.protocol_id,p.code)) AS protocol_ids
FROM protocol_eligibility_criteria pec
JOIN protocol_eligibility_criteria_values pecv ON pec.id = pecv.criteria_id
JOIN protocols p ON pecv.protocol_id = p.id
GROUP BY pec.id;

-- name: UnlinkEligibilityFromProtocol :exec
DELETE FROM protocol_eligibility_criteria_values
WHERE protocol_id = $1 AND criteria_id = $2;

-- name: VerifyLinkEligibilityToProtocol :one
SELECT * FROM protocol_eligibility_criteria_values
WHERE protocol_id = $1 AND criteria_id = $2;

-- name: GetEligibilityCriteriaByID :one
SELECT pec.*, ARRAY_AGG(ROW(pecv.protocol_id,p.code)) AS protocol_ids
FROM protocol_eligibility_criteria pec
JOIN protocol_eligibility_criteria_values pecv ON pec.id = pecv.criteria_id
JOIN protocols p ON pecv.protocol_id = p.id
WHERE pec.id = $1;

-- name: GetEligibilityCriteriaByType :many
SELECT pec.*, ARRAY_AGG(ROW(pecv.protocol_id,p.code)) AS protocol_ids
FROM protocol_eligibility_criteria pec
JOIN protocol_eligibility_criteria_values pecv ON pec.id = pecv.criteria_id
JOIN protocols p ON pecv.protocol_id = p.id
WHERE LOWER(pec.type) = LOWER($1)
GROUP BY pec.id;

-- name: GetEligibilityByProtocol :many
SELECT c.*
FROM protocol_eligibility_criteria c
JOIN protocol_eligibility_criteria_values v ON c.id = v.criteria_id
WHERE v.protocol_id = $1;
