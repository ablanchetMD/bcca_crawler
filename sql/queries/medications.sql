-- name: AddMedication :one
INSERT INTO medications (name, description, category,alternate_names)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: AddPrescription :one
INSERT INTO medication_prescription (medication_id, dose, route, frequency, duration, instructions, renewals)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpsertMedication :one
WITH input_values(id, name, description, category,alternate_names) AS (
  VALUES (
    CASE 
      WHEN @id = '00000000-0000-0000-0000-000000000000'::uuid 
      THEN gen_random_uuid() 
      ELSE @id 
    END,
    @name,
    @description,
    @category,
    @alternate_names::TEXT[]
  )
)
INSERT INTO medications (id, name, description, category,alternate_names)
SELECT id, name, description, category FROM input_values
ON CONFLICT (id) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    category = EXCLUDED.category,
    alternate_names = EXCLUDED.alternate_names,
    updated_at = NOW()
RETURNING *;

-- name: UpsertPrescription :one
WITH input_values(id, medication_id, dose, route, frequency, duration, instructions, renewals) AS (
  VALUES (
    CASE 
      WHEN @id = '00000000-0000-0000-0000-000000000000'::uuid 
      THEN gen_random_uuid() 
      ELSE @id::uuid 
    END,
    @medication_id::uuid,
    @dose,
    @route::prescription_route_enum,
    @frequency,
    @duration,
    @instructions,
    @renewals::int
  )
)
INSERT INTO medication_prescription (id, medication_id, dose, route, frequency, duration, instructions, renewals)
SELECT id, medication_id, dose, route, frequency, duration, instructions, renewals FROM input_values
ON CONFLICT (id) DO UPDATE
SET medication_id = EXCLUDED.medication_id,
    dose = EXCLUDED.dose,
    route = EXCLUDED.route,
    frequency = EXCLUDED.frequency,
    duration = EXCLUDED.duration,
    instructions = EXCLUDED.instructions,
    renewals = EXCLUDED.renewals,
    updated_at = NOW()
RETURNING *;

-- name: GetPrescriptions :many
SELECT m.id as medication_id, m.name, m.description, m.category,m.alternate_names, p.id as medication_prescription_id, p.dose, p.created_at,p.updated_at, p.route, p.frequency, p.duration, p.instructions, p.renewals
FROM medications m
JOIN medication_prescription p ON m.id = p.medication
ORDER BY m.name ASC;

-- name: GetPrescriptionsByMed :many
SELECT m.id as medication_id, m.name, m.description, m.category,m.alternate_names, p.id as medication_prescription_id, p.dose, p.created_at,p.updated_at, p.route, p.frequency, p.duration, p.instructions, p.renewals
FROM medications m
JOIN medication_prescription p ON m.id = p.medication
WHERE m.id = $1;

-- name: GetPrescriptionByID :one
SELECT m.id as medication_id, m.name, m.description, m.category,m.alternate_names, p.id as medication_prescription_id, p.dose,p.created_at,p.updated_at, p.route, p.frequency, p.duration, p.instructions, p.renewals
FROM medications m
JOIN medication_prescription p ON m.id = p.medication
WHERE p.id = $1;

-- name: RemovePrescription :exec
DELETE FROM medication_prescription
WHERE id = $1;

-- name: AddPrescriptionToProtocolByCategory :exec
INSERT INTO protocol_meds (protocol_id, prescription_id, category)
VALUES ($1, $2, $3)
ON CONFLICT DO NOTHING;

-- name: RemovePrescriptionFromProtocolByCategory :exec
DELETE FROM protocol_meds
WHERE protocol_id = $1 AND prescription_id = $2 AND category = $3;

-- name: GetPrescriptionsByProtocolByCategory :many
SELECT m.id as medication_id, m.name, m.description, m.category,m.alternate_names, p.id as medication_prescription_id, p.dose, p.created_at,p.updated_at, p.route, p.frequency, p.duration, p.instructions, p.renewals
FROM medications m
JOIN medication_prescription p ON m.id = p.medication
JOIN protocol_meds pm ON p.id = pm.prescription_id
WHERE pm.protocol_id = $1 AND pm.category = $2;

-- name: AddMedicationModification :one
INSERT INTO medication_modifications (category, subcategory, adjustment, medication_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: UpdateMedicationModification :one
UPDATE medication_modifications
SET
    updated_at = NOW(),
    category = $2,
    subcategory = $3,
    adjustment = $4,
    medication_id = $5
WHERE id = $1
RETURNING *;

-- name: UpsertMedicationModification :one
INSERT INTO medication_modifications (id, category, subcategory, adjustment, medication_id)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (id) DO UPDATE
SET category = EXCLUDED.category,
    subcategory = EXCLUDED.subcategory,
    adjustment = EXCLUDED.adjustment,
    medication_id = EXCLUDED.medication_id,
    updated_at = NOW()
RETURNING *;

-- name: RemoveMedicationModification :exec
DELETE FROM medication_modifications
WHERE id = $1;

-- name: GetModificationsByMedication :many
SELECT * FROM medication_modifications
WHERE medication_id = $1;

-- name: GetMedicationModificationsByProtocol :many
SELECT m.id as medication_id, m.name, m.description, m.category,m.alternate_names, mod.id as modification_id, mod.category as modification_category, mod.subcategory as modification_subcategory, mod.adjustment
FROM medication_modifications mod
JOIN medications m ON mod.medication_id = m.id
JOIN protocol_treatment pt ON m.id = pt.medication
JOIN treatment_cycles_values tc ON pt.id = tc.protocol_treatment_id
JOIN protocol_cycles pc ON tc.protocol_cycles_id = pc.id
WHERE pc.protocol_id = $1;

-- name: GetMedicationModificationsByMedication :many
SELECT m.id as medication_id, mod.id as modification_id, mod.category as modification_category, mod.subcategory as modification_subcategory, mod.adjustment
FROM medication_modifications mod
JOIN medications m ON mod.medication_id = m.id
WHERE m.id = $1;

-- name: GetMedicationModificationByID :one
SELECT * FROM medication_modifications
WHERE id = $1;

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

-- name: DeleteMedication :exec
DELETE FROM medications
WHERE id = $1;
