-- name: CreateProtocol :one
INSERT INTO protocols (id, created_at, updated_at, tumor_group, code, name, tags, notes)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: UpdateProtocol :one
UPDATE protocols
SET
    tumor_group = $2,
    updated_at = NOW(),
    code = $3,
    name = $4,
    tags = $5,
    notes = $6
WHERE id = $1
RETURNING *;

 --name: DeleteProtocol :exec
DELETE FROM protocols
WHERE id = $1;

-- name: GetProtocolByID :one
SELECT * FROM protocols
WHERE id = $1;

-- name: GetProtocols :many
SELECT * FROM protocols
ORDER BY name ASC;

-- name: GetProtocolsByTumorGroupAsc :many
SELECT * FROM protocols
WHERE tumor_group = $1
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: GetProtocolsByTumorGroupDesc :many
SELECT * FROM protocols
WHERE tumor_group = $1
ORDER BY name DESC
LIMIT $1 OFFSET $2;

-- name: GetProtocolsByTumorGroupAndTagsAsc :many
SELECT * FROM protocols
WHERE tumor_group = $1
AND tags @> $2
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: GetProtocolsByTumorGroupAndTagsDesc :many
SELECT * FROM protocols
WHERE tumor_group = $1
AND tags @> $2
ORDER BY name DESC
LIMIT $1 OFFSET $2;





