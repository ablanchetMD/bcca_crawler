-- name: CreateCancer :one
INSERT INTO cancers (code,name,tumor_group,tags,notes)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: UpdateCancer :one
UPDATE cancers
SET
    updated_at = NOW(),
    code = $2,
    name = $3,
    tumor_group = $4,
    tags = $5,
    notes = $6
WHERE id = $1
RETURNING *;

-- name: DeleteCancer :exec
DELETE FROM cancers
WHERE id = $1;

-- name: GetCancerByID :one
SELECT * FROM cancers
WHERE id = $1;

-- name: GetCancers :many
SELECT * FROM cancers
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: GetCancersByTags :many
SELECT * FROM cancers
WHERE tags @> $1
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: GetCancersOnlyTumorGroupAsc :many
SELECT * FROM cancers
WHERE tumor_group = $1
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: GetCancersOnlyTumorGroupAndTagsAsc :many
SELECT * FROM cancers
WHERE tumor_group = $1
AND tags @> $2
ORDER BY name ASC
LIMIT $3 OFFSET $4;

-- name: AddProtocolToCancer :exec
INSERT INTO cancer_protocols (cancer_id, protocol_id)
VALUES (
    $1,
    $2
);

-- name: RemoveProtocolFromCancer :exec
DELETE FROM cancer_protocols
WHERE cancer_id = $1 AND protocol_id = $2;

-- name: GetProtocolsForCancer :many
SELECT p.* 
FROM cancer_protocols cp
JOIN protocols p ON cp.protocol_id = p.id
WHERE cp.cancer_id = $1;
