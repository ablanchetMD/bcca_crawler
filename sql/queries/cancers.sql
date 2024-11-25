 -- name: CreateCancer :one
INSERT INTO cancers (code,name,tags,notes)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: UpdateCancer :one
UPDATE cancers
SET
    updated_at = NOW(),
    code = $2,
    name = $3,
    tags = $4,
    notes = $5
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
LIMIT $1 OFFSET $2;

-- name: GetCancersByTagsDesc :many
SELECT * FROM cancers
WHERE tags @> $1
ORDER BY name DESC
LIMIT $1 OFFSET $2;

