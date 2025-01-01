-- name: AddMedication :one
INSERT INTO medications (name, description, category)
VALUES ($1, $2, $3)
RETURNING *;

-- name: AddPrescription :one
INSERT INTO medication_prescription (medication, dose, route, frequency, duration, instructions, renewals)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: RemovePrescription :exec
DELETE FROM medication_prescription
WHERE id = $1;

-- name: AddPreMedicationToProtocol :exec
INSERT INTO protocol_pre_medications_values (protocol_id, medication_prescription_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: AddSupportiveMedicationToProtocol :exec
INSERT INTO protocol_supportive_medication_values (protocol_id, medication_prescription_id)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: RemovePreMedicationFromProtocol :exec
DELETE FROM protocol_pre_medications_values
WHERE protocol_id = $1 AND medication_prescription_id = $2;

-- name: RemoveSupportiveMedicationFromProtocol :exec
DELETE FROM protocol_supportive_medication_values
WHERE protocol_id = $1 AND medication_prescription_id = $2;

-- name: GetSupportiveMedicationsByProtocol :many
SELECT m.id as medication_id, m.name, m.description, m.category, p.id as medication_prescription_id, p.dose, p.route, p.frequency, p.duration, p.instructions, p.renewals
FROM medications m
JOIN medication_prescription p ON m.id = p.medication
JOIN protocol_supportive_medication_values s ON p.id = s.medication_prescription_id
WHERE s.protocol_id = $1;

-- name: GetPreMedicationsByProtocol :many
SELECT m.id as medication_id, m.name, m.description, m.category, p.id as medication_prescription_id, p.dose, p.route, p.frequency, p.duration, p.instructions, p.renewals
FROM medications m
JOIN medication_prescription p ON m.id = p.medication
JOIN protocol_pre_medications_values s ON p.id = s.medication_prescription_id
WHERE s.protocol_id = $1;

-- name: AddMedicationModification :one
INSERT INTO medication_modifications (category, description, adjustment, medication_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateMedicationModification :one
UPDATE medication_modifications
SET
    updated_at = NOW(),
    category = $2,
    description = $3,
    adjustment = $4,
    medication_id = $5
WHERE id = $1
RETURNING *;

-- name: RemoveMedicationModification :exec
DELETE FROM medication_modifications
WHERE id = $1;

-- name: GetModificationsByMedication :many
SELECT * FROM medication_modifications
WHERE medication_id = $1;

-- name: GetMedicationModificationsByProtocol :many
SELECT m.id as medication_id, m.name, m.description, m.category, mod.id as modification_id, mod.category as modification_category, mod.description as modification_description, mod.adjustment
FROM medication_modifications mod
JOIN medications m ON mod.medication_id = m.id
JOIN protocol_treatment pt ON m.id = pt.medication
JOIN treatment_cycles_values tc ON pt.id = tc.protocol_treatment_id
JOIN protocol_cycles pc ON tc.protocol_cycles_id = pc.id
WHERE pc.protocol_id = $1;

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
