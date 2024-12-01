-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (    
    $1,
    $2,
    $3,
    $4,
    $5    
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
RETURNING *;

-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password,role)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UpdateUser :one
UPDATE users
SET
    updated_at = $1,
    email = $2,
    password = $3,
    role = $4,
    is_verified = $5,
    deleted_at = $6,
    deleted_by = $7,
    last_active = $8
WHERE id = $9
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;


