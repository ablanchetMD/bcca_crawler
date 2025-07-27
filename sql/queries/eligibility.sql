-- name: InsertEligibilityCriteria :one
INSERT INTO protocol_eligibility_criteria (type, description)
VALUES ($1, $2)
RETURNING *;

-- name: UpsertEligibilityCriteria :one
WITH input_values(id, type, description) AS (
    VALUES
    (
        CASE
            WHEN @id = '00000000-0000-0000-0000-000000000000'::uuid 
            THEN gen_random_uuid() 
            ELSE @id 
        END,        
        @type::eligibility_enum,
        @description
    )
)
INSERT INTO protocol_eligibility_criteria (id, type, description)
SELECT id, type, description FROM input_values
ON CONFLICT (id) DO UPDATE
SET type = EXCLUDED.type::eligibility_enum,
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
SELECT 
    pec.*, 
    COALESCE(
        (
            SELECT json_agg(
            json_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
                -- 'created_at', p.created_at,
                -- 'updated_at', p.updated_at
            )
        )
        FROM protocol_eligibility_criteria_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.criteria_id = pec.id
        ),
        '[]'
    ) AS protocol_ids
FROM 
    protocol_eligibility_criteria pec;

-- name: GetEligibilityCriteriaByID :one
SELECT 
    pec.*, 
    COALESCE(
        (
            SELECT json_agg(
            json_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
                -- 'created_at', p.created_at,
                -- 'updated_at', p.updated_at
            )
        )
        FROM protocol_eligibility_criteria_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.criteria_id = pec.id
        ),
        '[]'
    ) AS protocol_ids
FROM 
    protocol_eligibility_criteria pec
WHERE
    pec.id = $1;

-- name: GetEligibilityCriteriaByType :many
SELECT 
    pec.*, 
    COALESCE(
        (
            SELECT json_agg(
            json_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
                -- 'created_at', p.created_at,
                -- 'updated_at', p.updated_at
            )
        )
        FROM protocol_eligibility_criteria_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.criteria_id = pec.id
        ),
        '[]'
    ) AS protocol_ids
FROM 
    protocol_eligibility_criteria pec
WHERE 
    LOWER(pec.type) = LOWER($1);

-- name: UnlinkEligibilityFromProtocol :exec
DELETE FROM protocol_eligibility_criteria_values
WHERE protocol_id = $1 AND criteria_id = $2;

-- name: VerifyLinkEligibilityToProtocol :one
SELECT * FROM protocol_eligibility_criteria_values
WHERE protocol_id = $1 AND criteria_id = $2;

-- name: GetEligibilityByProtocol :many
SELECT c.*
FROM protocol_eligibility_criteria c
JOIN protocol_eligibility_criteria_values v ON c.id = v.criteria_id
WHERE v.protocol_id = $1;

-- name: UpdateEligibilityProtocols :exec
WITH current_protocols AS (
    SELECT pcv.protocol_id 
    FROM protocol_eligibility_criteria_values pcv 
    WHERE pcv.criteria_id = $1
),
to_remove AS (
    DELETE FROM protocol_eligibility_criteria_values pcv
    WHERE pcv.criteria_id = $1
    AND pcv.protocol_id NOT IN (SELECT unnest(@protocol_ids::uuid[]))
    RETURNING pcv.protocol_id
),
to_add AS (
    INSERT INTO protocol_eligibility_criteria_values (criteria_id, protocol_id)
    SELECT $1, new_protocol
    FROM unnest(@protocol_ids::uuid[]) AS new_protocol
    WHERE new_protocol NOT IN (SELECT cp.protocol_id FROM current_protocols cp)
    RETURNING protocol_id
)
SELECT 
    (SELECT COUNT(*) FROM to_remove) AS removed, 
    (SELECT COUNT(*) FROM to_add) AS added;