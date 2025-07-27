
-- name: ResetDatabase :exec
TRUNCATE TABLE users, cancers, protocols, refresh_tokens, logs, physicians,
protocol_eligibility_criteria, protocol_cautions, medications,
medication_prescription, protocol_treatment, medication_modifications,
protocol_precautions, protocol_ppos, tests, protocol_tests, protocol_meds,
protocol_cycles, article_references, toxicities, toxicity_grades,
protocol_tox_modifications RESTART IDENTITY CASCADE;