-- name: AddMedication :one
INSERT INTO medications (name, description, category)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AddManyMedications :many
INSERT INTO medications (name, description, category)
VALUES ($1::TEXT[], $2::TEXT[], $3::TEXT[])
ON CONFLICT (name) DO NOTHING
RETURNING *;

-- name: AddPreMedication :one
INSERT INTO protocol_premedication (medication, dose, route, frequency, duration, notes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: AddSupportiveMedication :one
INSERT INTO protocol_supportive_medication (medication, dose, route, frequency, duration, notes)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: RemoveSupportiveMedication :exec
DELETE FROM protocol_supportive_medication
WHERE id = $1;

-- name: RemovePreMedication :exec
DELETE FROM protocol_premedication
WHERE id = $1;

-- name: AddPreMedicationToProtocol :exec
INSERT INTO protocol_pre_medications_values (protocol_id, pre_medication_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: AddSupportiveMedicationToProtocol :exec
INSERT INTO protocol_supportive_medication_values (protocol_id, supportive_medication_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemovePreMedicationFromProtocol :exec
DELETE FROM protocol_pre_medications_values
WHERE protocol_id = $1 AND pre_medication_id = $2;

-- name: RemoveSupportiveMedicationFromProtocol :exec
DELETE FROM protocol_supportive_medication_values
WHERE protocol_id = $1 AND supportive_medication_id = $2;

-- name: GetSupportiveMedicationsByProtocol :many
SELECT m.id as medication_id, m.name, m.description, m.category, p.id as supportive_medication_id, p.dose, p.route, p.frequency, p.duration, p.notes
FROM medications m
JOIN protocol_supportive_medication p ON m.id = p.medication
JOIN protocol_supportive_medication_values s ON p.id = s.supportive_medication_id
WHERE s.protocol_id = $1;


-- name: GetPreMedicationsByProtocol :many
SELECT m.id as medication_id, m.name, m.description, m.category, s.id as pre_medication_id, s.dose, s.route, s.frequency, s.duration, s.notes
FROM medications m
JOIN protocol_premedication s ON m.id = s.medication
JOIN protocol_pre_medications_values p ON s.id = p.pre_medication_id
WHERE p.protocol_id = $1;

-- name: GetMedications :many
SELECT * FROM medications
ORDER BY name ASC;

-- name: GetMedicationByID :one
SELECT * FROM medications
WHERE id = $1;

-- name: GetMedicationsByCategory :many
SELECT * FROM medications
WHERE category = $1
ORDER BY name ASC;

-- name: GetMedicationByName :one
SELECT * FROM medications
WHERE name = $1;
