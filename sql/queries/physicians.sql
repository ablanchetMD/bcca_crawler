-- name: CreatePhysician :one
INSERT INTO physicians (first_name, last_name, email, site)
VALUES (
    $1,
    $2,
    $3,
    $4
)    
RETURNING *;

-- name: UpdatePhysician :one
UPDATE physicians
SET
    updated_at = NOW(),
    first_name = $2,
    last_name = $3,
    email = $4,
    site = $5
WHERE id = $1
returning *;

-- name: UpsertPhysician :one
INSERT INTO physicians (id, first_name, last_name, email, site)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE
SET first_name = EXCLUDED.first_name,
    last_name = EXCLUDED.last_name,
    email = EXCLUDED.email,
    site = EXCLUDED.site,
    updated_at = NOW()
RETURNING *;

-- name: DeletePhysician :exec
DELETE FROM physicians
WHERE id = $1;

-- name: GetPhysicianByName :one
SELECT * FROM physicians
WHERE first_name = $1 AND last_name = $2;

-- name: GetPhysicianByID :one
SELECT * FROM physicians
WHERE id = $1;

-- name: GetPhysicians :many
SELECT * FROM physicians
ORDER BY last_name ASC;

-- name: GetPhysiciansBySite :many
SELECT * FROM physicians
WHERE site = $1
ORDER BY last_name ASC;

-- name: GetPhysicianByProtocol :many
SELECT p.*
FROM physicians p
JOIN protocol_contact_physicians ON p.id = protocol_contact_physicians.physician_id
WHERE protocol_contact_physicians.protocol_id = $1
ORDER BY p.last_name ASC;

-- name: AddPhysicianToProtocol :exec
INSERT INTO protocol_contact_physicians (protocol_id, physician_id)
VALUES ($1, $2);

-- name: AddManyPhysicianToProtocol :exec
INSERT INTO protocol_contact_physicians (protocol_id, physician_id)
VALUES ($1::UUID[], $2::UUID[])
ON CONFLICT DO NOTHING;

-- name: RemovePhysicianFromProtocol :exec
DELETE FROM protocol_contact_physicians
WHERE protocol_id = $1 AND physician_id = $2;

