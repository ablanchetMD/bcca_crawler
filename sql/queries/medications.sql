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
      ELSE @id::uuid
    END,
    @name::text,
    @description::text,
    @category,
    COALESCE(@alternate_names::TEXT[], '{}')
  )
)
INSERT INTO medications (id, name, description, category,alternate_names)
SELECT id, name, description, category,alternate_names FROM input_values
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
JOIN medication_prescription p ON m.id = p.medication_id
ORDER BY m.name ASC;

-- name: GetPrescriptionsByMed :many
SELECT m.id as medication_id, m.name, m.description, m.category,m.alternate_names, p.id as medication_prescription_id, p.dose, p.created_at,p.updated_at, p.route, p.frequency, p.duration, p.instructions, p.renewals
FROM medications m
JOIN medication_prescription p ON m.id = p.medication_id
WHERE m.id = $1;

-- name: GetPrescriptionByID :one
SELECT m.id as medication_id, m.name, m.description, m.category,m.alternate_names, p.id as medication_prescription_id, p.dose,p.created_at,p.updated_at, p.route, p.frequency, p.duration, p.instructions, p.renewals
FROM medications m
JOIN medication_prescription p ON m.id = p.medication_id
WHERE p.id = $1;

-- name: GetPrescriptionsByArguments :one
SELECT p.id 
FROM medication_prescription p
WHERE p.medication_id = $1 AND p.dose = $2 AND p.route = $3
  AND p.frequency = $4 AND p.duration = $5 AND p.instructions = $6;

-- name: RemovePrescription :exec
DELETE FROM medication_prescription
WHERE id = $1;

-- name: UpsertProtoMedCategory :one
WITH input_values(id, protocol_id, category, comments) AS (
  VALUES (
    CASE 
      WHEN @id::uuid  = '00000000-0000-0000-0000-000000000000'
      THEN gen_random_uuid() 
      ELSE @id::uuid 
    END,
    @protocol_id::uuid,
    @category::TEXT,
    @comments::text   
  )
)
INSERT INTO protocol_meds (id, protocol_id, category, comments)
SELECT id, protocol_id, category, comments FROM input_values
ON CONFLICT (id) DO UPDATE
SET protocol_id = EXCLUDED.protocol_id,
    category = EXCLUDED.category,
    comments = EXCLUDED.comments,
    updated_at = NOW()
RETURNING *;

-- name: GetMedCategoryByID :one
SELECT jsonb_build_object(
  'id', pm.id,
  'created_at', pm.created_at,
  'updated_at', pm.updated_at,
  'category', pm.category,
  'comments', pm.comments,
  'medications', COALESCE((
    SELECT jsonb_agg(
      jsonb_build_object(            
        'id', mp.id,
        'dose', mp.dose,
        'route', mp.route,
        'frequency', mp.frequency,
        'duration', mp.duration,
        'instructions', mp.instructions,
        'renewals', mp.renewals,
        'medication_id', m.id,
        'medication_name', m.name,
        'medication_description', m.description,
        'medication_category', m.category,
        'medication_alternates', m.alternate_names
      )
      ORDER BY m.name
    )
    FROM protocol_meds_values pmv
    JOIN medication_prescription mp ON mp.id = pmv.medication_prescription_id
    JOIN medications m ON m.id = mp.medication_id
    WHERE pmv.protocol_meds_id = pm.id
  ), '[]'::jsonb)
)
FROM protocol_meds pm
WHERE pm.id = $1;

-- name: RemoveMedCategoryByID :exec
DELETE FROM protocol_meds
WHERE id = $1;

-- name: AddPrescriptionToProtocolCategory :exec
INSERT INTO protocol_meds_values (protocol_meds_id, medication_prescription_id)
VALUES ($1, $2);

-- name: RemovePrescriptionFromProtocolCategory :exec
DELETE FROM protocol_meds_values
WHERE protocol_meds_id = $1 AND medication_prescription_id = $2;

-- name: GetPrescriptionsByProtocolByCategory :many
SELECT m.id as medication_id, m.name, m.description, m.category,m.alternate_names, p.id as medication_prescription_id, p.dose, p.created_at,p.updated_at, p.route, p.frequency, p.duration, p.instructions, p.renewals
FROM medications m
JOIN medication_prescription p ON m.id = p.medication_id
JOIN protocol_meds pm ON p.id = pm.prescription_id
WHERE pm.protocol_id = $1 AND pm.category = $2;


-- name: GetProtocolPrescriptions :one
SELECT COALESCE(jsonb_agg(protocol_prescriptions ORDER BY category ASC), '[]'::jsonb) AS data
FROM (
  SELECT 
    pro.id,
    pro.created_at,
    pro.updated_at,
    pro.category,
    pro.comments,     
    COALESCE((
      SELECT jsonb_agg(prescription_data)
      FROM (
        SELECT 
          jsonb_build_object(            
            'id', mp.id,
            'dose', mp.dose,
            'route', mp.route,
            'frequency', mp.frequency,
            'duration', mp.duration,
            'instructions', mp.instructions,
            'renewals', mp.renewals,
            'medication_id', m.id,
            'medication_name', m.name,
            'medication_description', m.description,
            'medication_category',m.category,
            'medication_alternates',m.alternate_names
          ) AS prescription_data                 
        FROM protocol_meds_values pmv
        JOIN medication_prescription mp ON mp.id = pmv.medication_prescription_id
        JOIN medications m ON m.id = mp.medication_id
        WHERE pmv.protocol_meds_id = pro.id
        ORDER BY m.name
      ) t
    ), '[]'::jsonb) AS medications
  FROM protocol_meds pro  
  WHERE pro.protocol_id = $1
) AS protocol_prescriptions;

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
JOIN protocol_treatment pt ON m.id = pt.medication_id
JOIN treatment_cycles_values tc ON pt.id = tc.protocol_treatment_id
JOIN protocol_cycles pc ON tc.protocol_cycles_id = pc.id
WHERE pc.protocol_id = $1;

-- name: GetProtocolMedicationsWithModifications :many
SELECT 
  m.id AS medication_id,
  m.name AS medication_name,
  jsonb_agg(
    jsonb_build_object(
      'category', grouped.category,
      'subcategories', grouped.subcategories
    )
  ) AS categories
FROM protocol_treatment pt
JOIN medications m ON m.id = pt.medication_id

JOIN LATERAL (
  SELECT 
    mm1.category,
    jsonb_agg(
      jsonb_build_object(
        'subcategory', mm1.subcategory,
        'adjustment', mm1.adjustment
      )
    ) AS subcategories
  FROM medication_modifications mm1
  WHERE mm1.medication_id = m.id
  GROUP BY mm1.category
) AS grouped ON true

WHERE pt.id IN (
  SELECT ptv.protocol_treatment_id
  FROM treatment_cycles_values ptv
  JOIN protocol_cycles pc ON pc.id = ptv.protocol_cycles_id
  WHERE pc.protocol_id = $1
)

GROUP BY m.id, m.name;
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
