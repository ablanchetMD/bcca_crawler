-- name: InsertEligibilityCriteria :one
INSERT INTO protocol_eligibility_criteria (type, description)
VALUES ($1, $2)
RETURNING id, created_at, updated_at, type, description;

-- name: InsertManyEligibilityCriteria :many
INSERT INTO protocol_eligibility_criteria (type, description)
VALUES ($1::TEXT[], $2::TEXT[])
ON CONFLICT DO NOTHING
RETURNING id, created_at, updated_at, type, description;

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


-- name: LinkEligibilityToProtocol :exec
INSERT INTO protocol_eligibility_criteria_values (protocol_id, criteria_id)
VALUES ($1::UUID[], $2::UUID[])
ON CONFLICT DO NOTHING;

-- name: UnlinkEligibilityFromProtocol :exec
DELETE FROM protocol_eligibility_criteria_values
WHERE protocol_id = $1 AND criteria_id = $2;

-- name: GetEligibilityCriteriaByID :one
SELECT * FROM protocol_eligibility_criteria
WHERE id = $1;

-- name: GetEligibilityCriteriaBy :many
SELECT * FROM protocol_eligibility_criteria
WHERE type = $1;

-- name: GetEligibilityByProtocol :many
SELECT c.id, c.type, c.description
FROM protocol_eligibility_criteria c
JOIN protocol_eligibility_criteria_values v ON c.id = v.criteria_id
WHERE v.protocol_id = $1;
