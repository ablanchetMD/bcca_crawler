-- name: AddProtocolTreatment :one
INSERT INTO protocol_treatment (medication, dose, route, frequency, duration, administration_guide)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProtocolTreatment :one
UPDATE protocol_treatment
SET
    updated_at = NOW(),
    medication = $2,
    dose = $3,
    route = $4,
    frequency = $5,
    duration = $6,
    administration_guide = $7
WHERE id = $1
RETURNING *;

-- name: GetProtocolTreatmentByData :one
SELECT * FROM protocol_treatment
WHERE medication = $1 AND dose = $2 AND route = $3 AND frequency = $4 AND duration = $5 AND administration_guide = $6;


-- name: AddTreatmentModification :one
INSERT INTO treatment_modifications (category, description, adjustement, treatment_id)
VALUES ($1, $2, $3, $4)
RETURNING *;


-- name: UpdateTreatmentModification :one
UPDATE treatment_modifications
SET
    updated_at = NOW(),
    category = $2,
    description = $3,
    adjustement = $4,
    treatment_id = $5
WHERE id = $1
RETURNING *;

-- name: RemoveTreatmentModification :exec
DELETE FROM treatment_modifications
WHERE id = $1;

-- name: AddToxicityModification :one
INSERT INTO toxicity_modifications (title, grade, adjustement, protocol_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateToxicityModification :one
UPDATE toxicity_modifications
SET
    updated_at = NOW(),
    title = $2,
    grade = $3,
    adjustement = $4,
    protocol_id = $5
WHERE id = $1
RETURNING *;

-- name: RemoveToxicityModification :exec
DELETE FROM toxicity_modifications
WHERE id = $1;

-- name: AddCycleToProtocol :one
INSERT INTO protocol_cycles (protocol_id, cycle, cycle_duration)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AddTreatmentToCycle :exec
INSERT INTO treatment_cycles_junction (protocol_cycles_id, protocol_treatment_id)
VALUES ($1, $2);

-- name: RemoveTreatmentFromCycle :exec
DELETE FROM treatment_cycles_junction
WHERE protocol_cycles_id = $1 AND protocol_treatment_id = $2;

-- name: RemoveProtocolTreatment :exec
DELETE FROM protocol_treatment
WHERE id = $1;

-- name: GetProtocolTreatmentByID :one
SELECT * FROM protocol_treatment
WHERE id = $1;

-- name: GetCyclesByProtocol :many
SELECT protocol_cycles.*
FROM protocol_cycles
WHERE protocol_cycles.protocol_id = $1
ORDER BY protocol_cycles.cycle ASC;

-- name: GetTreatmentsByCycle :many
SELECT protocol_treatment.*
FROM protocol_treatment
JOIN treatment_cycles_junction ON protocol_treatment.id = treatment_cycles_junction.protocol_treatment_id
WHERE treatment_cycles_junction.protocol_cycles_id = $1
ORDER BY protocol_treatment.medication ASC;

-- name: GetTreatmentModificationsByTreatment :many
SELECT * FROM treatment_modifications
WHERE treatment_id = $1;

-- name: GetToxicityModificationsByProtocol :many
SELECT * FROM toxicity_modifications
WHERE protocol_id = $1;



