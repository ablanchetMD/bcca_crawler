-- name: CreateProtocolPrecaution :one
INSERT INTO protocol_precautions (title, description)
VALUES ($1, $2)    
RETURNING *;

-- name: GetProtocolPrecautionByID :one
SELECT * FROM protocol_precautions WHERE id = $1;

-- name: GetProtocolPrecautionByTitleAndDescription :one
SELECT * FROM protocol_precautions WHERE title = $1 AND description = $2;

-- name: UpdateProtocolPrecaution :one
UPDATE protocol_precautions SET title = $2, description = $3 WHERE id = $1 RETURNING *;

-- name: DeleteProtocolPrecaution :exec
DELETE FROM protocol_precautions WHERE id = $1;

-- name: AddProtocolPrecautionToProtocol :exec
INSERT INTO protocol_precautions_values (protocol_id, precaution_id) VALUES ($1, $2);

-- name: RemoveProtocolPrecautionFromProtocol :exec
DELETE FROM protocol_precautions_values WHERE protocol_id = $1 AND precaution_id = $2;

-- name: GetProtocolPrecautionsByProtocol :many
SELECT p.* FROM protocol_precautions p JOIN protocol_precautions_values v ON p.id = v.precaution_id WHERE v.protocol_id = $1;

-- name: CreateProtocolCaution :one
INSERT INTO protocol_cautions (description)
VALUES ($1) 
RETURNING *;

-- name: UpsertCaution :one
WITH input_values(id, description) AS (
  VALUES (
    CASE 
      WHEN $1 = '00000000-0000-0000-0000-000000000000'::uuid 
      THEN gen_random_uuid() 
      ELSE $1 
    END,
    $2
  )
)
INSERT INTO protocol_cautions (id, description)
SELECT id, description FROM input_values
ON CONFLICT (id) DO UPDATE
SET description = EXCLUDED.description,
    updated_at = NOW()
RETURNING *;

-- name: GetCautionWithProtocols :many
SELECT 
    pec.*, 
    COALESCE(
        (
            SELECT json_agg(
            json_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
            )
        )
        FROM protocol_cautions_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.caution_id = pec.id
        ),
        '[]'
    ) AS protocol_ids
FROM 
    protocol_cautions pec;

-- name: GetCautionByIDWithProtocols :one
SELECT 
    pec.*, 
    COALESCE(
        (
            SELECT json_agg(
            json_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
            )
        )
        FROM protocol_cautions_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.caution_id = pec.id
        ),
        '[]'
    ) AS protocol_ids
FROM 
    protocol_cautions pec
WHERE
    pec.id = $1;


-- name: UpdateCautionProtocols :exec
WITH current_protocols AS (
    SELECT pcv.protocol_id 
    FROM protocol_cautions_values pcv 
    WHERE pcv.caution_id = $1
),
to_remove AS (
    DELETE FROM protocol_cautions_values pcv
    WHERE pcv.caution_id = $1
    AND pcv.protocol_id NOT IN (SELECT unnest($2::uuid[]))
    RETURNING pcv.protocol_id
),
to_add AS (
    INSERT INTO protocol_cautions_values (caution_id, protocol_id)
    SELECT $1, new_protocol
    FROM unnest($2::uuid[]) AS new_protocol
    WHERE new_protocol NOT IN (SELECT cp.protocol_id FROM current_protocols cp)
    RETURNING protocol_id
)
SELECT 
    (SELECT COUNT(*) FROM to_remove) AS removed, 
    (SELECT COUNT(*) FROM to_add) AS added;

-- name: UpsertPrecaution :one
WITH input_values(id,title, description) AS (
  VALUES (
    CASE 
      WHEN $1 = '00000000-0000-0000-0000-000000000000'::uuid 
      THEN gen_random_uuid() 
      ELSE $1 
    END,
    $2,
    $3
  )
)
INSERT INTO protocol_precautions (id,title, description)
SELECT id,title, description FROM input_values
ON CONFLICT (id) DO UPDATE
SET description = EXCLUDED.description,
    title = EXCLUDED.title,
    updated_at = NOW()
RETURNING *;

-- name: GetPrecautionWithProtocols :many
SELECT 
    pec.*, 
    COALESCE(
        (
            SELECT json_agg(
            json_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
            )
        )
        FROM protocol_precautions_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.precaution_id = pec.id
        ),
        '[]'
    ) AS protocol_ids
FROM 
    protocol_precautions pec;


-- name: GetPrecautionByIDWithProtocols :one
SELECT 
    pec.*, 
    COALESCE(
        (
            SELECT json_agg(
            json_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
            )
        )
        FROM protocol_precautions_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.precaution_id = pec.id
        ),
        '[]'
    ) AS protocol_ids
FROM 
    protocol_precautions pec
WHERE
    pec.id = $1;

-- name: UpdatePrecautionProtocols :exec
WITH current_protocols AS (
    SELECT pcv.protocol_id 
    FROM protocol_precautions_values pcv 
    WHERE pcv.precaution_id = @precaution_id
),
to_remove AS (
    DELETE FROM protocol_precautions_values pcv
    WHERE pcv.precaution_id = @precaution_id
    AND pcv.protocol_id NOT IN (SELECT unnest(@protocol_ids::uuid[]))
    RETURNING pcv.protocol_id
),
to_add AS (
    INSERT INTO protocol_precautions_values (precaution_id, protocol_id)
    SELECT @precaution_id, new_protocol
    FROM unnest(@protocol_ids::uuid[]) AS new_protocol
    WHERE new_protocol NOT IN (SELECT cp.protocol_id FROM current_protocols cp)
    RETURNING protocol_id
)
SELECT 
    (SELECT COUNT(*) FROM to_remove) AS removed, 
    (SELECT COUNT(*) FROM to_add) AS added;

-- name: GetProtocolCautionByID :one
SELECT * FROM protocol_cautions WHERE id = $1;

-- name: GetProtocolCautionByDescription :one
SELECT * FROM protocol_cautions WHERE description = $1;

-- name: UpdateProtocolCaution :one
UPDATE protocol_cautions SET description = $2 WHERE id = $1 RETURNING *;

-- name: DeleteProtocolCaution :exec
DELETE FROM protocol_cautions WHERE id = $1;

-- name: AddProtocolCautionToProtocol :exec
INSERT INTO protocol_cautions_values (protocol_id, caution_id) VALUES ($1, $2);

-- name: RemoveProtocolCautionFromProtocol :exec
DELETE FROM protocol_cautions_values WHERE protocol_id = $1 AND caution_id = $2;

-- name: GetProtocolCautionsByProtocol :many
SELECT c.* FROM protocol_cautions c JOIN protocol_cautions_values v ON c.id = v.caution_id WHERE v.protocol_id = $1;
