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
WITH tox_grades AS (
    SELECT
        tox.id AS toxicity_id,
        COALESCE(
            json_agg(DISTINCT json_build_object(
            'modification_id', ptm.id,
            'toxicity_grade_id', tg.id,
            'grade', tg.grade,
            'description', tg.description,
            'created_at', ptm.created_at,
            'updated_at', ptm.updated_at,
            'adjustment', ptm.adjustment
        ) ORDER BY tg.grade) FILTER (WHERE ptm.id IS NOT NULL), '[]') AS grades
    FROM
        toxicities tox
        JOIN toxicity_grades tg ON tg.toxicity_id = tox.id
        JOIN protocol_tox_modifications ptm ON ptm.toxicity_grade_id = tg.id
    GROUP BY tox.id
),
categorized_tests AS (
  SELECT
    pt.protocol_id,
    json_object_agg(
      t.test_category, 
      json_object_agg(
        t.urgency,
        COALESCE(
        json_agg(DISTINCT json_build_object(
          'test_id', t.id,
          'test_name', t.name,
          'description', t.description,
          'form_url', t.form_url,
          'unit', t.unit,
          'lower_limit', t.lower_limit,
          'upper_limit', t.upper_limit
        )) FILTER (WHERE t.id IS NOT NULL), '[]'
        )
      )
    ) AS categorized_tests
  FROM protocol_tests pt
  JOIN tests t ON pt.test_id = t.id
  GROUP BY pt.protocol_id
),
categorized_meds AS (
  SELECT
    pm.protocol_id,
    json_object_agg(
      pm.category,
      COALESCE(json_agg(DISTINCT json_build_object(
        'medication_id', m.id,
        'medication_name', m.name,
        'dose', mp.dose,
        'route', mp.route,
        'frequency', mp.frequency,
        'duration', mp.duration,
        'instructions', mp.instructions
      ))FILTER (WHERE m.id IS NOT NULL), '[]'
    )) AS categorized_meds
  FROM protocol_meds pm
  JOIN medication_prescription mp ON pm.prescription_id = mp.id
  JOIN medications m ON mp.medication = m.id
  GROUP BY pm.protocol_id
)
SELECT
    p.id,
    p.tumor_group,
    p.code,
    p.name,
    p.tags,
    p.created_at,
    p.updated_at,
    p.revised_on,
    p.activated_on,
    -- Aggregate associated physicians
    COALESCE(json_agg(DISTINCT json_build_object(
        'physician_id', ph.id,
        'first_name', ph.first_name,
        'last_name', ph.last_name,
        'email', ph.email,
        'site', ph.site
    )) FILTER (WHERE ph.id IS NOT NULL), '[]'::json) AS associated_physicians,
    -- Aggregate associated medications
    COALESCE(ct.categorized_tests, '[]') AS categorized_tests,
    COALESCE(cm.categorized_meds, '[]') AS categorized_medications,
    -- Aggregate eligibility criteria
    COALESCE(json_agg(DISTINCT json_build_object(
        'criteria_id', pec.id,
        'type', pec.type,
        'description', pec.description
    )) FILTER (WHERE pec.id IS NOT NULL), '[]'::json) AS eligibility_criteria,
    -- Aggregate cautions
    COALESCE(json_agg(DISTINCT json_build_object(
        'caution_id', pc.id,
        'description', pc.description
    )) FILTER (WHERE pc.id IS NOT NULL), '[]'::json) AS protocol_cautions,
    -- Aggregate precautions
    COALESCE(json_agg(DISTINCT json_build_object(
        'precaution_id', pp.id,
        'title', pp.title,
        'description', pp.description
    )) FILTER (WHERE pp.id IS NOT NULL), '[]'::json) AS protocol_precautions,
    -- Aggregate cycles
    COALESCE(json_agg(DISTINCT json_build_object(
        'cycle_id', pc.id,
        'cycle', pc.cycle,
        'cycle_duration', pc.cycle_duration
    )) FILTER (WHERE pc.id IS NOT NULL), '[]'::json) AS protocol_cycles,
    -- Aggregate references
    COALESCE(json_agg(DISTINCT json_build_object(
        'reference_id', ar.id,
        'title', ar.title,
        'authors', ar.authors,
        'journal', ar.journal,
        'year', ar.year,
        'pmid', ar.pmid,
        'doi', ar.doi
    )) FILTER (WHERE ar.id IS NOT NULL), '[]'::json) AS article_references,
    -- Aggregate toxicity modifications
    COALESCE(
        json_agg(
            DISTINCT json_build_object(
            'id', tox.id,
            'created_at',tox.created_at,
            'updated_at',tox.updated_at,
            'title',tox.title,
            'category',tox.category,
            'description',tox.description,
            'grades_with_adjustment', tgrades.grades                 
    )) FILTER (WHERE tox.id IS NOT NULL), '[]'::json) AS toxicities_with_adjustments
FROM
    protocols p
    -- Join with physicians associated with the protocol
    LEFT JOIN protocol_contact_physicians pcp ON p.id = pcp.protocol_id
    LEFT JOIN physicians ph ON pcp.physician_id = ph.id      
    -- Join with eligibility criteria
    LEFT JOIN protocol_eligibility_criteria_values pecv ON p.id = pecv.protocol_id
    LEFT JOIN protocol_eligibility_criteria pec ON pecv.criteria_id = pec.id
    -- Join with cautions
    LEFT JOIN protocol_cautions_values pcv ON p.id = pcv.protocol_id
    LEFT JOIN protocol_cautions pc ON pcv.caution_id = pc.id
    -- Join with precautions
    LEFT JOIN protocol_precautions_values ppv ON p.id = ppv.protocol_id
    LEFT JOIN protocol_precautions pp ON ppv.precaution_id = pp.id
    -- Join with cycles
    LEFT JOIN protocol_cycles pc ON p.id = pc.protocol_id
    -- Join with references
    LEFT JOIN protocol_references_value prv ON p.id = prv.protocol_id
    LEFT JOIN article_references ar ON prv.reference_id = ar.id
    -- Join with toxicity modifications
    LEFT JOIN protocol_tox_modifications ptm ON p.id = ptm.protocol_id
    LEFT JOIN toxicity_grades tg ON ptm.toxicity_grade_id = tg.id
    LEFT JOIN toxicities tox ON tg.toxicity_id = tox.id
    LEFT JOIN tox_grades tgrades ON tox.id = tgrades.toxicity_id
    --meds and tests
    LEFT JOIN categorized_tests ct ON p.id = ct.protocol_id
    LEFT JOIN categorized_meds cm ON p.id = cm.protocol_id
WHERE
    p.id = $1    
GROUP BY
    p.id;

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
