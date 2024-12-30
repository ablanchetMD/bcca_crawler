-- name: InsertEligibilityCriteria :one
INSERT INTO protocol_eligibility_criteria (type, description)
VALUES ($1, $2)
RETURNING *;

-- name: InsertManyEligibilityCriteria :exec
INSERT INTO protocol_eligibility_criteria (type, description)
VALUES ($1::UUID[], $2::UUID[])
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

-- name: LinkManyEligibilityToProtocol :exec
INSERT INTO protocol_eligibility_criteria_values (protocol_id, criteria_id)
VALUES ($1::UUID[], $2::UUID[])
ON CONFLICT DO NOTHING;

-- name: LinkEligibilityToProtocol :exec
INSERT INTO protocol_eligibility_criteria_values (protocol_id, criteria_id)
VALUES ($1, $2);


-- name: UnlinkEligibilityFromProtocol :exec
DELETE FROM protocol_eligibility_criteria_values
WHERE protocol_id = $1 AND criteria_id = $2;

-- name: GetEligibilityCriteriaByID :one
SELECT * FROM protocol_eligibility_criteria
WHERE id = $1;

-- name: GetEligibilityCriteriaByType :many
SELECT * FROM protocol_eligibility_criteria
WHERE type = $1;

-- name: GetEligibilityByProtocol :many
SELECT c.*
FROM protocol_eligibility_criteria c
JOIN protocol_eligibility_criteria_values v ON c.id = v.criteria_id
WHERE v.protocol_id = $1;
