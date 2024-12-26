-- name: CreatePhysician :one
INSERT INTO physicians (first_name, last_name, email, site)
VALUES (
    $1,
    $2,
    $3,
    $4
)
returning *;

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

-- name: DeletePhysician :exec
DELETE FROM physicians
WHERE id = $1;

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
SELECT physicians.id, physicians.first_name, physicians.last_name, physicians.email, physicians.site
FROM physicians
JOIN protocol_contact_physicians ON physicians.id = protocol_contact_physicians.physician_id
WHERE protocol_contact_physicians.protocol_id = $1
ORDER BY last_name ASC;

-- name: AddPhysicianToProtocol :exec
INSERT INTO protocol_contact_physicians (protocol_id, physician_id)
VALUES ($1::UUID[], $2::UUID[]);

-- name: RemovePhysicianFromProtocol :exec
DELETE FROM protocol_contact_physicians
WHERE protocol_id = $1 AND physician_id = $2;

