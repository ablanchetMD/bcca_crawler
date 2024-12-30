// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: physicians.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const addManyPhysicianToProtocol = `-- name: AddManyPhysicianToProtocol :exec
INSERT INTO protocol_contact_physicians (protocol_id, physician_id)
VALUES ($1::UUID[], $2::UUID[])
ON CONFLICT DO NOTHING
`

type AddManyPhysicianToProtocolParams struct {
	Column1 []uuid.UUID
	Column2 []uuid.UUID
}

func (q *Queries) AddManyPhysicianToProtocol(ctx context.Context, arg AddManyPhysicianToProtocolParams) error {
	_, err := q.db.ExecContext(ctx, addManyPhysicianToProtocol, pq.Array(arg.Column1), pq.Array(arg.Column2))
	return err
}

const addPhysicianToProtocol = `-- name: AddPhysicianToProtocol :exec
INSERT INTO protocol_contact_physicians (protocol_id, physician_id)
VALUES ($1, $2)
`

type AddPhysicianToProtocolParams struct {
	ProtocolID  uuid.UUID
	PhysicianID uuid.UUID
}

func (q *Queries) AddPhysicianToProtocol(ctx context.Context, arg AddPhysicianToProtocolParams) error {
	_, err := q.db.ExecContext(ctx, addPhysicianToProtocol, arg.ProtocolID, arg.PhysicianID)
	return err
}

const createPhysician = `-- name: CreatePhysician :one
INSERT INTO physicians (first_name, last_name, email, site)
VALUES (
    $1,
    $2,
    $3,
    $4
)    
RETURNING id, created_at, updated_at, first_name, last_name, email, site
`

type CreatePhysicianParams struct {
	FirstName string
	LastName  string
	Email     string
	Site      string
}

func (q *Queries) CreatePhysician(ctx context.Context, arg CreatePhysicianParams) (Physician, error) {
	row := q.db.QueryRowContext(ctx, createPhysician,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Site,
	)
	var i Physician
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Site,
	)
	return i, err
}

const deletePhysician = `-- name: DeletePhysician :exec
DELETE FROM physicians
WHERE id = $1
`

func (q *Queries) DeletePhysician(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deletePhysician, id)
	return err
}

const getPhysicianByID = `-- name: GetPhysicianByID :one
SELECT id, created_at, updated_at, first_name, last_name, email, site FROM physicians
WHERE id = $1
`

func (q *Queries) GetPhysicianByID(ctx context.Context, id uuid.UUID) (Physician, error) {
	row := q.db.QueryRowContext(ctx, getPhysicianByID, id)
	var i Physician
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Site,
	)
	return i, err
}

const getPhysicianByName = `-- name: GetPhysicianByName :one
SELECT id, created_at, updated_at, first_name, last_name, email, site FROM physicians
WHERE first_name = $1 AND last_name = $2
`

type GetPhysicianByNameParams struct {
	FirstName string
	LastName  string
}

func (q *Queries) GetPhysicianByName(ctx context.Context, arg GetPhysicianByNameParams) (Physician, error) {
	row := q.db.QueryRowContext(ctx, getPhysicianByName, arg.FirstName, arg.LastName)
	var i Physician
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Site,
	)
	return i, err
}

const getPhysicianByProtocol = `-- name: GetPhysicianByProtocol :many
SELECT p.id, p.created_at, p.updated_at, p.first_name, p.last_name, p.email, p.site
FROM physicians p
JOIN protocol_contact_physicians ON p.id = protocol_contact_physicians.physician_id
WHERE protocol_contact_physicians.protocol_id = $1
ORDER BY p.last_name ASC
`

func (q *Queries) GetPhysicianByProtocol(ctx context.Context, protocolID uuid.UUID) ([]Physician, error) {
	rows, err := q.db.QueryContext(ctx, getPhysicianByProtocol, protocolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Physician
	for rows.Next() {
		var i Physician
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Site,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPhysicians = `-- name: GetPhysicians :many
SELECT id, created_at, updated_at, first_name, last_name, email, site FROM physicians
ORDER BY last_name ASC
`

func (q *Queries) GetPhysicians(ctx context.Context) ([]Physician, error) {
	rows, err := q.db.QueryContext(ctx, getPhysicians)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Physician
	for rows.Next() {
		var i Physician
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Site,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPhysiciansBySite = `-- name: GetPhysiciansBySite :many
SELECT id, created_at, updated_at, first_name, last_name, email, site FROM physicians
WHERE site = $1
ORDER BY last_name ASC
`

func (q *Queries) GetPhysiciansBySite(ctx context.Context, site string) ([]Physician, error) {
	rows, err := q.db.QueryContext(ctx, getPhysiciansBySite, site)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Physician
	for rows.Next() {
		var i Physician
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Site,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removePhysicianFromProtocol = `-- name: RemovePhysicianFromProtocol :exec
DELETE FROM protocol_contact_physicians
WHERE protocol_id = $1 AND physician_id = $2
`

type RemovePhysicianFromProtocolParams struct {
	ProtocolID  uuid.UUID
	PhysicianID uuid.UUID
}

func (q *Queries) RemovePhysicianFromProtocol(ctx context.Context, arg RemovePhysicianFromProtocolParams) error {
	_, err := q.db.ExecContext(ctx, removePhysicianFromProtocol, arg.ProtocolID, arg.PhysicianID)
	return err
}

const updatePhysician = `-- name: UpdatePhysician :one
UPDATE physicians
SET
    updated_at = NOW(),
    first_name = $2,
    last_name = $3,
    email = $4,
    site = $5
WHERE id = $1
returning id, created_at, updated_at, first_name, last_name, email, site
`

type UpdatePhysicianParams struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Email     string
	Site      string
}

func (q *Queries) UpdatePhysician(ctx context.Context, arg UpdatePhysicianParams) (Physician, error) {
	row := q.db.QueryRowContext(ctx, updatePhysician,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.Email,
		arg.Site,
	)
	var i Physician
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Site,
	)
	return i, err
}
