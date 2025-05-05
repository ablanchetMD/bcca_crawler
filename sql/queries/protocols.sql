-- name: CreateProtocol :one
INSERT INTO protocols (tumor_group, code, name, tags, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateProtocolbyScraping :one
INSERT INTO protocols (tumor_group, code, name, tags, notes, revised_on, activated_on)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateProtocol :one
UPDATE protocols
SET
    tumor_group = $2,
    updated_at = NOW(),
    code = $3,
    name = $4,
    tags = $5,
    notes = $6,
    protocol_url = $7,
    patient_handout_url = $8,
    revised_on = $9,
    activated_on = $10
WHERE id = $1
RETURNING *;

-- name: UpsertProtocol :one
WITH input_values(id, tumor_group, code,name,tags,notes,protocol_url,patient_handout_url,revised_on,activated_on) AS (
    VALUES
    (
        CASE
            WHEN $1 = '00000000-0000-0000-0000-000000000000'::uuid 
            THEN gen_random_uuid() 
            ELSE $1 
        END,        
        $2::tumor_group_enum,
        $3,
        $4,
        $5::TEXT[],
        $6,
        $7,
        $8,
        $9,
        $10        
    )
)
INSERT INTO protocols (id, tumor_group, code, name, tags, notes, protocol_url, patient_handout_url, revised_on, activated_on)
SELECT id, tumor_group,code,name,tags,notes,protocol_url,patient_handout_url,revised_on,activated_on FROM input_values
ON CONFLICT (id) DO UPDATE
SET tumor_group = EXCLUDED.tumor_group::tumor_group_enum,
    code = EXCLUDED.code,
    name = EXCLUDED.name,
    tags = EXCLUDED.tags,
    notes = EXCLUDED.notes,
    protocol_url = EXCLUDED.protocol_url,
    patient_handout_url = EXCLUDED.patient_handout_url,
    revised_on = EXCLUDED.revised_on,
    activated_on = EXCLUDED.activated_on,    
    updated_at = NOW()
RETURNING *;

-- name: DeleteProtocol :exec
DELETE FROM protocols
WHERE id = $1;

-- name: GetProtocolByID :one
SELECT * FROM protocols
WHERE id = $1;

-- name: GetProtocolByCode :one
SELECT * FROM protocols
WHERE code = $1;

-- name: GetProtocolsAsc :many
SELECT * FROM protocols
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: GetProtocolsDesc :many
SELECT * FROM protocols
ORDER BY name DESC
LIMIT $1 OFFSET $2;

-- name: GetProtocolsOnlyTumorGroupAsc :many
SELECT * FROM protocols
WHERE tumor_group = $1
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: GetProtocolsOnlyTumorGroupDesc :many
SELECT * FROM protocols
WHERE tumor_group = $1
ORDER BY name DESC
LIMIT $2 OFFSET $3;

-- name: GetProtocolsOnlyTumorGroupAndTagsAsc :many
SELECT * FROM protocols
WHERE tumor_group = $1
AND tags @> $2
ORDER BY name ASC
LIMIT $3 OFFSET $4;

-- name: GetProtocolsOnlyTumorGroupAndTagsDesc :many
SELECT * FROM protocols
WHERE tumor_group = $1
AND tags @> $2
ORDER BY name DESC
LIMIT $3 OFFSET $4;
