-- name: CreateProtocol :one
INSERT INTO protocols (tumor_group, code, name, tags, notes)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CreateProtocolbyScraping :one
INSERT INTO protocols (tumor_group, code, name, tags, notes, revised_on, activated_on)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateProtocol :one
UPDATE protocols
SET
    tumor_group = $2,
    updated_at = NOW(),
    code = $3,
    name = $4,
    tags = $5,
    notes = $6,
    protocol_url = $7,
    patient_handout_url = $8,
    revised_on = $9,
    activated_on = $10
WHERE id = $1
RETURNING *;

-- name: UpsertProtocol :one
WITH input_values(id, tumor_group, code,name,tags,notes,protocol_url,patient_handout_url,revised_on,activated_on) AS (
    VALUES
    (
        CASE
            WHEN @id = '00000000-0000-0000-0000-000000000000'::uuid 
            THEN gen_random_uuid() 
            ELSE @id 
        END,        
        @tumor_group::tumor_group_enum,
        @code,
        @name,
        @tags::TEXT[],
        @notes,
        @protocol_url,
        @patient_handout_url,
        @revised_on,
        @activated_on       
    )
)
INSERT INTO protocols (id, tumor_group, code, name, tags, notes, protocol_url, patient_handout_url, revised_on, activated_on)
SELECT id, tumor_group,code,name,tags,notes,protocol_url,patient_handout_url,revised_on,activated_on FROM input_values
ON CONFLICT (id) DO UPDATE
SET tumor_group = EXCLUDED.tumor_group::tumor_group_enum,
    code = EXCLUDED.code,
    name = EXCLUDED.name,
    tags = EXCLUDED.tags,
    notes = EXCLUDED.notes,
    protocol_url = EXCLUDED.protocol_url,
    patient_handout_url = EXCLUDED.patient_handout_url,
    revised_on = EXCLUDED.revised_on,
    activated_on = EXCLUDED.activated_on,    
    updated_at = NOW()
RETURNING *;

-- name: GetProtocolData :one
SELECT sqlc.embed(protocols),sqlc.embed(protocol_cycles),
FROM protocols

JOIN article_references ON article_references.id = protocol_references_value.reference_id
LEFT JOIN protocol_references_value ON protocol_references_value.protocol_id = protocols.id

JOIN physicians ON physicians.id = protocol_contact_physicians.physician_id
LEFT JOIN protocol_contact_physicians ON protocol_contact_physicians.protocol_id = protocols.id

JOIN protocol_eligibility_criteria ON protocol_eligibility_criteria.id = protocol_eligibility_criteria_values.criteria_id
LEFT JOIN protocol_eligibility_criteria_values ON protocol_eligibility_criteria_values.protocol_id = protocols.id

JOIN protocol_cautions ON protocol_cautions.id = protocol_cautions_values.caution_id
LEFT JOIN protocol_cautions_values ON protocol_cautions_values.protocol_id = protocols.id

JOIN protocol_precautions ON protocol_precautions.id = protocol_precautions_values.precaution_id
LEFT JOIN protocol_precautions_values ON protocol_precautions_values.protocol_id = protocols.id

--cycles + treatment
JOIN protocol_cycles ON protocol_cycles.protocol_id = protocols.id
JOIN protocol_treatment ON protocol_treatment.id = treatment_cycles_values.protocol_treatment_id
LEFT JOIN treatment_cycles_values ON treatment_cycles_values.protocol_cycles_id = protocol_cycles.id

--toxicities + totixicies mod + grades
JOIN toxicities ON toxicities.id = toxicity_grades.toxicity_id
JOIN toxicity_grades ON toxicity_grades.id = protocol_tox_modifications.toxicity_grade_id
LEFT JOIN protocol_tox_modifications ON protocol_tox_modifications.protocol_id = protocols.id

--tests
JOIN tests ON tests.id = protocol_tests.test_id
LEFT JOIN protocol_tests ON protocol_tests.protocol_id = protocols.id

--meds
JOIN medications ON medications.id = medication_prescription.medication
JOIN medication_prescription ON medication_prescription.id = protocol_meds.prescription_id
LEFT JOIN protocol_meds ON protocol_meds.protocol_id = protocols.id

WHERE protocols.id = $1

-- name: DeleteProtocol :exec
DELETE FROM protocols
WHERE id = $1;

-- name: GetProtocolByID :one
SELECT * FROM protocols
WHERE id = $1;

-- name: GetProtocolByCode :one
SELECT * FROM protocols
WHERE code = $1;

-- name: GetProtocolsAsc :many
SELECT * FROM protocols
ORDER BY name ASC
LIMIT $1 OFFSET $2;

-- name: GetProtocolsDesc :many
SELECT * FROM protocols
ORDER BY name DESC
LIMIT $1 OFFSET $2;

-- name: GetProtocolsOnlyTumorGroupAsc :many
SELECT * FROM protocols
WHERE tumor_group = $1
ORDER BY name ASC
LIMIT $2 OFFSET $3;

-- name: GetProtocolsOnlyTumorGroupDesc :many
SELECT * FROM protocols
WHERE tumor_group = $1
ORDER BY name DESC
LIMIT $2 OFFSET $3;

-- name: GetProtocolsOnlyTumorGroupAndTagsAsc :many
SELECT * FROM protocols
WHERE tumor_group = $1
AND tags @> $2
ORDER BY name ASC
LIMIT $3 OFFSET $4;

-- name: GetProtocolsOnlyTumorGroupAndTagsDesc :many
SELECT * FROM protocols
WHERE tumor_group = $1
AND tags @> $2
ORDER BY name DESC
LIMIT $3 OFFSET $4;
