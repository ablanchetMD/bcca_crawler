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

-- name: AddManyProtocolPrecautionToProtocol :exec
INSERT INTO protocol_precautions_values (protocol_id, precaution_id) VALUES ($1::UUID[], $2::UUID[]) ON CONFLICT DO NOTHING;

-- name: CreateProtocolCaution :one
INSERT INTO protocol_cautions (description)
VALUES ($1) 
RETURNING *;

-- name: UpsertCaution :one
INSERT INTO protocol_cautions (id, description)
VALUES ($1, $2)
ON CONFLICT (id) DO UPDATE
SET description = EXCLUDED.description,
    updated_at = NOW()
RETURNING *;

-- name: GetCautionWithProtocols :many
SELECT pec.*, ARRAY_AGG(ROW(pecv.protocol_id,p.code)) AS protocol_ids
FROM protocol_cautions pec
JOIN protocol_cautions_values pecv ON pec.id = pecv.caution_id
JOIN protocols p ON pecv.protocol_id = p.id
GROUP BY pec.id;

-- name: GetCautionByIDWithProtocols :one
SELECT pec.*, ARRAY_AGG(ROW(pecv.protocol_id,p.code)) AS protocol_ids
FROM protocol_cautions pec
JOIN protocol_cautions_values pecv ON pec.id = pecv.caution_id
JOIN protocols p ON pecv.protocol_id = p.id
WHERE pec.id = $1;

-- name: UpsertPrecaution :one
INSERT INTO protocol_precautions (id,title, description)
VALUES ($1, $2, $3)
ON CONFLICT (id) DO UPDATE
SET description = EXCLUDED.description,
    title = EXCLUDED.title,
    updated_at = NOW()
RETURNING *;

-- name: GetPrecautionWithProtocols :many
SELECT pec.*, ARRAY_AGG(ROW(pecv.protocol_id,p.code)) AS protocol_ids
FROM protocol_precautions pec
JOIN protocol_precautions_values pecv ON pec.id = pecv.precaution_id
JOIN protocols p ON pecv.protocol_id = p.id
GROUP BY pec.id;

-- name: GetPrecautionByIDWithProtocols :one
SELECT pec.*, ARRAY_AGG(ROW(pecv.protocol_id,p.code)) AS protocol_ids
FROM protocol_precautions pec
JOIN protocol_precautions_values pecv ON pec.id = pecv.precaution_id
JOIN protocols p ON pecv.protocol_id = p.id
WHERE pec.id = $1;

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

-- name: AddManyProtocolCautionToProtocol :exec
INSERT INTO protocol_cautions_values (protocol_id, caution_id) VALUES ($1::UUID[], $2::UUID[]) ON CONFLICT DO NOTHING;