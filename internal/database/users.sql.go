// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (    
    $1,
    $2,
    $3,
    $4,
    $5    
)
RETURNING token, created_at, updated_at, expires_at, revoked_at, user_id
`

type CreateRefreshTokenParams struct {
	Token     string
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	ExpiresAt time.Time
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken,
		arg.Token,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.ExpiresAt,
	)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
		&i.UserID,
	)
	return i, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, password,role)
VALUES (
    gen_random_uuid(),
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, created_at, updated_at, email, password, role, is_verified, deleted_at, deleted_by, last_active
`

type CreateUserParams struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	Email     string
	Password  string
	Role      string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Email,
		arg.Password,
		arg.Role,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.IsVerified,
		&i.DeletedAt,
		&i.DeletedBy,
		&i.LastActive,
	)
	return i, err
}

const deleteUsers = `-- name: DeleteUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteUsers)
	return err
}

const getRefreshToken = `-- name: GetRefreshToken :one
SELECT token, created_at, updated_at, expires_at, revoked_at, user_id FROM refresh_tokens WHERE token = $1
`

func (q *Queries) GetRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, getRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
		&i.UserID,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, created_at, updated_at, email, password, role, is_verified, deleted_at, deleted_by, last_active FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.IsVerified,
		&i.DeletedAt,
		&i.DeletedBy,
		&i.LastActive,
	)
	return i, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
RETURNING token, created_at, updated_at, expires_at, revoked_at, user_id
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, revokeRefreshToken, token)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExpiresAt,
		&i.RevokedAt,
		&i.UserID,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
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
RETURNING id, created_at, updated_at, email, password, role, is_verified, deleted_at, deleted_by, last_active
`

type UpdateUserParams struct {
	UpdatedAt  time.Time
	Email      string
	Password   string
	Role       string
	IsVerified bool
	DeletedAt  sql.NullTime
	DeletedBy  uuid.NullUUID
	LastActive sql.NullTime
	ID         uuid.UUID
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.UpdatedAt,
		arg.Email,
		arg.Password,
		arg.Role,
		arg.IsVerified,
		arg.DeletedAt,
		arg.DeletedBy,
		arg.LastActive,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Email,
		&i.Password,
		&i.Role,
		&i.IsVerified,
		&i.DeletedAt,
		&i.DeletedBy,
		&i.LastActive,
	)
	return i, err
}
