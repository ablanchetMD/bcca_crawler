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

-- name: UpsertProtocolTreatment :one
INSERT INTO protocol_treatment (id, medication, dose, route, frequency, duration, administration_guide)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (id) DO UPDATE
SET
    updated_at = NOW(),
    medication = EXCLUDED.medication,
    dose = EXCLUDED.dose,
    route = EXCLUDED.route,
    frequency = EXCLUDED.frequency,
    duration = EXCLUDED.duration,
    administration_guide = EXCLUDED.administration_guide
RETURNING *;

-- name: GetTreatments :many
SELECT m.id as medication_id, m.name as medication_name, m.description as medication_description, m.category as medication_category ,m.alternate_names as medication_alternates, pt.id as treatment_id, pt.dose, pt.created_at,pt.updated_at, pt.route, pt.frequency, pt.duration, pt.administration_guide
FROM medications m
JOIN protocol_treatment pt ON m.id = pt.medication
ORDER BY medication_name ASC;

-- name: GetProtocolTreatmentByID :one
SELECT m.id as medication_id, m.name as medication_name, m.description as medication_description, m.category as medication_category ,m.alternate_names as medication_alternates, pt.id as treatment_id, pt.dose, pt.created_at,pt.updated_at, pt.route, pt.frequency, pt.duration, pt.administration_guide
FROM medications m
JOIN protocol_treatment pt ON m.id = pt.medication
WHERE pt.id = $1;

-- name: GetTreatmentsByCycle :many
SELECT m.id as medication_id, m.name as medication_name, m.description as medication_description, m.category as medication_category ,m.alternate_names as medication_alternates, pt.id as treatment_id, pt.dose, pt.created_at,pt.updated_at, pt.route, pt.frequency, pt.duration, pt.administration_guide
FROM medications m
JOIN protocol_treatment pt ON m.id = pt.medication
JOIN treatment_cycles_values ON pt.id = treatment_cycles_values.protocol_treatment_id
WHERE treatment_cycles_values.protocol_cycles_id = $1
ORDER BY medication_name ASC;

-- name: GetProtocolTreatmentByData :one
SELECT * FROM protocol_treatment
WHERE medication = $1 AND dose = $2 AND route = $3 AND frequency = $4 AND duration = $5;

-- name: AddCycleToProtocol :one
INSERT INTO protocol_cycles (protocol_id, cycle, cycle_duration)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpsertCycleToProtocol :one
INSERT INTO protocol_cycles (id, protocol_id, cycle, cycle_duration)
VALUES ($1, $2, $3, $4)
ON CONFLICT (id) DO UPDATE
SET
    updated_at = NOW(),
    protocol_id = EXCLUDED.protocol_id,
    cycle = EXCLUDED.cycle,
    cycle_duration = EXCLUDED.cycle_duration
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



-- name: GetCycles :many
SELECT * FROM protocol_cycles
ORDER BY cycle ASC;

-- name: GetCycleByID :one
SELECT * FROM protocol_cycles
WHERE id = $1;

-- name: RemoveCycleByID :exec
DELETE FROM protocol_cycles
WHERE id = $1;

-- name: GetCyclesByProtocol :many
SELECT protocol_cycles.*
FROM protocol_cycles
WHERE protocol_cycles.protocol_id = $1
ORDER BY protocol_cycles.cycle ASC;

