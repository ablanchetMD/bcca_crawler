-- name: AddProtocolTreatment :one
INSERT INTO protocol_treatment (medication_id, dose, route, frequency, duration, administration_guide)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProtocolTreatment :one
UPDATE protocol_treatment
SET
    updated_at = NOW(),
    medication_id = $2,
    dose = $3,
    route = $4,
    frequency = $5,
    duration = $6,
    administration_guide = $7
WHERE id = $1
RETURNING *;

-- name: UpsertProtocolTreatment :one
WITH input_values(id, medication_id, dose, route, frequency, duration, administration_guide) AS (
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
    @administration_guide
  )
)
INSERT INTO protocol_treatment (id, medication_id, dose, route, frequency, duration, administration_guide)
SELECT id, medication_id, dose, route, frequency, duration, administration_guide FROM input_values
ON CONFLICT (id) DO UPDATE
SET medication_id = EXCLUDED.medication_id,
    dose = EXCLUDED.dose,
    route = EXCLUDED.route,
    frequency = EXCLUDED.frequency,
    duration = EXCLUDED.duration,
    administration_guide = EXCLUDED.administration_guide,    
    updated_at = NOW()
RETURNING *;

-- name: GetTreatments :many
SELECT m.id as medication_id, m.name as medication_name, m.description as medication_description, m.category as medication_category ,m.alternate_names as medication_alternates, pt.id as id, pt.dose, pt.created_at,pt.updated_at, pt.route, pt.frequency, pt.duration, pt.administration_guide
FROM medications m
JOIN protocol_treatment pt ON m.id = pt.medication_id
ORDER BY medication_name ASC;

-- name: GetProtocolTreatmentByID :one
SELECT m.id as medication_id, m.name as medication_name, m.description as medication_description, m.category as medication_category ,m.alternate_names as medication_alternates, pt.id as id, pt.dose, pt.created_at,pt.updated_at, pt.route, pt.frequency, pt.duration, pt.administration_guide
FROM medications m
JOIN protocol_treatment pt ON m.id = pt.medication_id
WHERE pt.id = $1;

-- name: GetTreatmentsByCycle :many
SELECT m.id as medication_id, m.name as medication_name, m.description as medication_description, m.category as medication_category ,m.alternate_names as medication_alternates, pt.id as id, pt.dose, pt.created_at,pt.updated_at, pt.route, pt.frequency, pt.duration, pt.administration_guide
FROM medications m
JOIN protocol_treatment pt ON m.id = pt.medication_id
JOIN treatment_cycles_values ON pt.id = treatment_cycles_values.protocol_treatment_id
WHERE treatment_cycles_values.protocol_cycles_id = $1
ORDER BY medication_name ASC;

-- name: GetProtocolTreatmentByData :one
SELECT * FROM protocol_treatment
WHERE medication_id = $1 AND dose = $2 AND route = $3 AND frequency = $4 AND duration = $5;

-- name: AddCycleToProtocol :one
INSERT INTO protocol_cycles (protocol_id, cycle, cycle_duration)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpsertCycleToProtocol :one
WITH input_values(id, protocol_id, cycle, cycle_duration) AS (
  VALUES (
    CASE 
      WHEN @id::uuid  = '00000000-0000-0000-0000-000000000000'
      THEN gen_random_uuid() 
      ELSE @id::uuid 
    END,
    @protocol_id::uuid,
    @cycle::TEXT,
    @cycle_duration::text    
  )
)
INSERT INTO protocol_cycles (id, protocol_id, cycle, cycle_duration)
SELECT id, protocol_id, cycle, cycle_duration FROM input_values
ON CONFLICT (id) DO UPDATE
SET protocol_id = EXCLUDED.protocol_id,
    cycle = EXCLUDED.cycle,
    cycle_duration = EXCLUDED.cycle_duration,      
    updated_at = NOW()
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

-- name: GetProtocolCyclesWithTreatments :one
SELECT COALESCE(jsonb_agg(cycle_data ORDER BY cycle_order), '[]'::jsonb) AS data
FROM (
  SELECT 
    pc.id,
    pc.created_at,
    pc.updated_at,
    pc.cycle,
    pc.cycle_duration,
    COALESCE(NULLIF(regexp_replace(pc.cycle, '\D', '', 'g'), '')::int, 0) AS cycle_order,
    COALESCE((
      SELECT jsonb_agg(treatment_data ORDER BY treatment_order)
      FROM (
        SELECT 
          jsonb_build_object(
            'medication_id', m.id,
            'medication_name', m.name,
            'medication_description', m.description,
            'medication_category', m.category,
            'medication_alternates', m.alternate_names,
            'id', pt.id,
            'dose', pt.dose,
            'created_at', pt.created_at,
            'updated_at', pt.updated_at,
            'route', pt.route,
            'frequency', pt.frequency,
            'duration', pt.duration,
            'administration_guide', pt.administration_guide
          ) AS treatment_data,
          COALESCE(NULLIF(regexp_replace(pt.frequency, '\D', '', 'g'), '')::int, 0) AS treatment_order
        FROM protocol_treatment pt
        LEFT JOIN medications m ON pt.medication_id = m.id
        LEFT JOIN treatment_cycles_values tc 
          ON tc.protocol_treatment_id = pt.id AND tc.protocol_cycles_id = pc.id
        WHERE tc.protocol_cycles_id = pc.id
      ) t
    ), '[]'::jsonb) AS treatments
  FROM protocol_cycles pc
  WHERE pc.protocol_id = $1
) cycle_data;

-- name: RemoveCycleByID :exec
DELETE FROM protocol_cycles
WHERE id = $1;
