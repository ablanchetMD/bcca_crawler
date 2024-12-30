-- name: AddTest :one
INSERT INTO tests (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: AddManyTests :exec
INSERT INTO tests (name, description)
VALUES ($1::TEXT[], $2::TEXT[])
ON CONFLICT (name) DO NOTHING;

-- name: AddBaselineTest :exec
INSERT INTO protocol_baseline_tests (protocol_id, test_id)
VALUES ($1, $2);

-- name: AddNonUrgentTest :exec
INSERT INTO protocol_baseline_tests_non_urgent (protocol_id, test_id)
VALUES ($1, $2);

-- name: AddIfNecessaryTest :exec
INSERT INTO protocol_baseline_tests_if_necessary (protocol_id, test_id)
VALUES ($1, $2);

-- name: AddFollowupTest :exec
INSERT INTO protocol_followup_tests (protocol_id, test_id)
VALUES ($1, $2);

-- name: AddFollowupIfNecessaryTest :exec
INSERT INTO protocol_followup_tests_if_necessary (protocol_id, test_id)
VALUES ($1, $2);

-- name: RemoveBaselineTest :exec
DELETE FROM protocol_baseline_tests
WHERE protocol_id = $1 AND test_id = $2;

-- name: RemoveNonUrgentTest :exec
DELETE FROM protocol_baseline_tests_non_urgent
WHERE protocol_id = $1 AND test_id = $2;

-- name: RemoveIfNecessaryTest :exec
DELETE FROM protocol_baseline_tests_if_necessary
WHERE protocol_id = $1 AND test_id = $2;

-- name: RemoveFollowupTest :exec
DELETE FROM protocol_followup_tests
WHERE protocol_id = $1 AND test_id = $2;

-- name: RemoveFollowupIfNecessaryTest :exec
DELETE FROM protocol_followup_tests_if_necessary
WHERE protocol_id = $1 AND test_id = $2;

-- name: GetTestByName :one
SELECT * FROM tests WHERE name = $1;

-- name: GetBaselineTestsByProtocol :many
SELECT t.*
FROM tests t
JOIN protocol_baseline_tests pb ON t.id = pb.test_id
WHERE pb.protocol_id = $1;

-- name: GetNonUrgentTestsByProtocol :many
SELECT t.*
FROM tests t
JOIN protocol_baseline_tests_non_urgent pb ON t.id = pb.test_id
WHERE pb.protocol_id = $1;

-- name: GetIfNecessaryTestsByProtocol :many
SELECT t.*
FROM tests t
JOIN protocol_baseline_tests_if_necessary pb ON t.id = pb.test_id
WHERE pb.protocol_id = $1;

-- name: GetFollowupTestsByProtocol :many
SELECT t.*
FROM tests t
JOIN protocol_followup_tests pb ON t.id = pb.test_id
WHERE pb.protocol_id = $1;

-- name: GetFollowupIfNecessaryTestsByProtocol :many
SELECT t.*
FROM tests t
JOIN protocol_followup_tests_if_necessary pb ON t.id = pb.test_id
WHERE pb.protocol_id = $1;


