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
WHERE medication = $1 AND dose = $2 AND route = $3 AND frequency = $4 AND duration = $5;

-- name: AddCycleToProtocol :one
INSERT INTO protocol_cycles (protocol_id, cycle, cycle_duration)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetCycleByData :one
SELECT * FROM protocol_cycles
WHERE protocol_id = $1 AND cycle = $2 AND cycle_duration = $3;

-- name: AddTreatmentToCycle :exec
INSERT INTO treatment_cycles_values (protocol_cycles_id, protocol_treatment_id)
VALUES ($1, $2);

-- name: RemoveTreatmentFromCycle :exec
DELETE FROM treatment_cycles_values
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
JOIN treatment_cycles_values ON protocol_treatment.id = treatment_cycles_values.protocol_treatment_id
WHERE treatment_cycles_values.protocol_cycles_id = $1
ORDER BY protocol_treatment.medication ASC;