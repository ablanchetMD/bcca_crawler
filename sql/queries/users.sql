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
SELECT rt.*, u.role
FROM refresh_tokens rt
INNER JOIN users u ON rt.user_id = u.id
WHERE rt.token = $1;

-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
RETURNING *;

-- name: RevokeRefreshTokenByUserId :many
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE user_id = $1
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

-- name: GetUsers :many
SELECT * FROM users
ORDER BY email ASC
LIMIT $1 OFFSET $2;

-- name: GetUsersByRole :many
SELECT * FROM users
WHERE role = $1
ORDER BY email ASC
LIMIT $2 OFFSET $3;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserRoleByID :one
SELECT role FROM users WHERE id = $1;

-- name: DeleteUserByID :exec
DELETE FROM users WHERE id = $1;

-- name: UpdateUser :one
UPDATE users
SET
    updated_at = NOW(),
    email = COALESCE($2, email),
    password = COALESCE($3, password),
    role = COALESCE($4, role),
    is_verified = COALESCE($5, is_verified),
    deleted_at = COALESCE($6, deleted_at),
    deleted_by = COALESCE($7, deleted_by),
    last_active = COALESCE($8, last_active)
WHERE id = $1
    AND (
        $2 IS DISTINCT FROM email
        OR $3 IS DISTINCT FROM password
        OR $4 IS DISTINCT FROM role
        OR $5 IS DISTINCT FROM is_verified
        OR $6 IS DISTINCT FROM deleted_at
        OR $7 IS DISTINCT FROM deleted_by
        OR $8 IS DISTINCT FROM last_active        
    )
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;


