// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: protocol_cautions.sql

package database

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const addProtocolCautionToProtocol = `-- name: AddProtocolCautionToProtocol :exec
INSERT INTO protocol_cautions_values (protocol_id, caution_id) VALUES ($1, $2)
`

type AddProtocolCautionToProtocolParams struct {
	ProtocolID uuid.UUID `json:"protocol_id"`
	CautionID  uuid.UUID `json:"caution_id"`
}

func (q *Queries) AddProtocolCautionToProtocol(ctx context.Context, arg AddProtocolCautionToProtocolParams) error {
	_, err := q.db.ExecContext(ctx, addProtocolCautionToProtocol, arg.ProtocolID, arg.CautionID)
	return err
}

const addProtocolPrecautionToProtocol = `-- name: AddProtocolPrecautionToProtocol :exec
INSERT INTO protocol_precautions_values (protocol_id, precaution_id) VALUES ($1, $2)
`

type AddProtocolPrecautionToProtocolParams struct {
	ProtocolID   uuid.UUID `json:"protocol_id"`
	PrecautionID uuid.UUID `json:"precaution_id"`
}

func (q *Queries) AddProtocolPrecautionToProtocol(ctx context.Context, arg AddProtocolPrecautionToProtocolParams) error {
	_, err := q.db.ExecContext(ctx, addProtocolPrecautionToProtocol, arg.ProtocolID, arg.PrecautionID)
	return err
}

const createProtocolCaution = `-- name: CreateProtocolCaution :one
INSERT INTO protocol_cautions (description)
VALUES ($1) 
RETURNING id, created_at, updated_at, description
`

func (q *Queries) CreateProtocolCaution(ctx context.Context, description string) (ProtocolCaution, error) {
	row := q.db.QueryRowContext(ctx, createProtocolCaution, description)
	var i ProtocolCaution
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
	)
	return i, err
}

const createProtocolPrecaution = `-- name: CreateProtocolPrecaution :one
INSERT INTO protocol_precautions (title, description)
VALUES ($1, $2)    
RETURNING id, created_at, updated_at, title, description
`

type CreateProtocolPrecautionParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (q *Queries) CreateProtocolPrecaution(ctx context.Context, arg CreateProtocolPrecautionParams) (ProtocolPrecaution, error) {
	row := q.db.QueryRowContext(ctx, createProtocolPrecaution, arg.Title, arg.Description)
	var i ProtocolPrecaution
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
	)
	return i, err
}

const deleteProtocolCaution = `-- name: DeleteProtocolCaution :exec
DELETE FROM protocol_cautions WHERE id = $1
`

func (q *Queries) DeleteProtocolCaution(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteProtocolCaution, id)
	return err
}

const deleteProtocolPrecaution = `-- name: DeleteProtocolPrecaution :exec
DELETE FROM protocol_precautions WHERE id = $1
`

func (q *Queries) DeleteProtocolPrecaution(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteProtocolPrecaution, id)
	return err
}

const getCautionByIDWithProtocols = `-- name: GetCautionByIDWithProtocols :one
SELECT 
    pec.id, pec.created_at, pec.updated_at, pec.description, 
    COALESCE(
        (
            SELECT jsonb_agg(
            jsonb_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
                -- 'created_at', p.created_at,
                -- 'updated_at', p.updated_at
            )
        )
        FROM protocol_cautions_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.caution_id = pec.id
        ),
        '[]'
    )::jsonb AS protocol_ids
FROM 
    protocol_cautions pec
WHERE
    pec.id = $1
`

type GetCautionByIDWithProtocolsRow struct {
	ID          uuid.UUID       `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Description string          `json:"description"`
	ProtocolIds json.RawMessage `json:"protocol_ids"`
}

func (q *Queries) GetCautionByIDWithProtocols(ctx context.Context, id uuid.UUID) (GetCautionByIDWithProtocolsRow, error) {
	row := q.db.QueryRowContext(ctx, getCautionByIDWithProtocols, id)
	var i GetCautionByIDWithProtocolsRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
		&i.ProtocolIds,
	)
	return i, err
}

const getCautionWithProtocols = `-- name: GetCautionWithProtocols :many
SELECT 
    pec.id, pec.created_at, pec.updated_at, pec.description, 
    COALESCE(
        (
            SELECT jsonb_agg(
            jsonb_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
                -- 'created_at', p.created_at,
                -- 'updated_at', p.updated_at
            )
        )
        FROM protocol_cautions_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.caution_id = pec.id
        ),
        '[]'
    )::jsonb AS protocol_ids
FROM 
    protocol_cautions pec
`

type GetCautionWithProtocolsRow struct {
	ID          uuid.UUID       `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Description string          `json:"description"`
	ProtocolIds json.RawMessage `json:"protocol_ids"`
}

func (q *Queries) GetCautionWithProtocols(ctx context.Context) ([]GetCautionWithProtocolsRow, error) {
	rows, err := q.db.QueryContext(ctx, getCautionWithProtocols)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCautionWithProtocolsRow{}
	for rows.Next() {
		var i GetCautionWithProtocolsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Description,
			&i.ProtocolIds,
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

const getPrecautionByIDWithProtocols = `-- name: GetPrecautionByIDWithProtocols :one
SELECT 
    pec.id, pec.created_at, pec.updated_at, pec.title, pec.description, 
    COALESCE(
        (
            SELECT jsonb_agg(
            jsonb_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
                -- 'created_at', p.created_at,
                -- 'updated_at', p.updated_at
            )
        )
        FROM protocol_precautions_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.precaution_id = pec.id
        ),
        '[]'
    )::jsonb AS protocol_ids
FROM 
    protocol_precautions pec
WHERE
    pec.id = $1
`

type GetPrecautionByIDWithProtocolsRow struct {
	ID          uuid.UUID       `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	ProtocolIds json.RawMessage `json:"protocol_ids"`
}

func (q *Queries) GetPrecautionByIDWithProtocols(ctx context.Context, id uuid.UUID) (GetPrecautionByIDWithProtocolsRow, error) {
	row := q.db.QueryRowContext(ctx, getPrecautionByIDWithProtocols, id)
	var i GetPrecautionByIDWithProtocolsRow
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
		&i.ProtocolIds,
	)
	return i, err
}

const getPrecautionWithProtocols = `-- name: GetPrecautionWithProtocols :many
SELECT 
    pec.id, pec.created_at, pec.updated_at, pec.title, pec.description, 
    COALESCE(
        (
            SELECT jsonb_agg(
            jsonb_build_object(
                'id', pecv.protocol_id, 
                'code', p.code
                -- 'created_at', p.created_at,
                -- 'updated_at', p.updated_at
            )
        )
        FROM protocol_precautions_values pecv
        JOIN protocols p ON pecv.protocol_id = p.id
        WHERE pecv.precaution_id = pec.id
        ),
        '[]'
    )::jsonb AS protocol_ids
FROM 
    protocol_precautions pec
`

type GetPrecautionWithProtocolsRow struct {
	ID          uuid.UUID       `json:"id"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	ProtocolIds json.RawMessage `json:"protocol_ids"`
}

func (q *Queries) GetPrecautionWithProtocols(ctx context.Context) ([]GetPrecautionWithProtocolsRow, error) {
	rows, err := q.db.QueryContext(ctx, getPrecautionWithProtocols)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetPrecautionWithProtocolsRow{}
	for rows.Next() {
		var i GetPrecautionWithProtocolsRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Description,
			&i.ProtocolIds,
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

const getProtocolCautionByDescription = `-- name: GetProtocolCautionByDescription :one
SELECT id, created_at, updated_at, description FROM protocol_cautions WHERE description = $1
`

func (q *Queries) GetProtocolCautionByDescription(ctx context.Context, description string) (ProtocolCaution, error) {
	row := q.db.QueryRowContext(ctx, getProtocolCautionByDescription, description)
	var i ProtocolCaution
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
	)
	return i, err
}

const getProtocolCautionByID = `-- name: GetProtocolCautionByID :one
SELECT id, created_at, updated_at, description FROM protocol_cautions WHERE id = $1
`

func (q *Queries) GetProtocolCautionByID(ctx context.Context, id uuid.UUID) (ProtocolCaution, error) {
	row := q.db.QueryRowContext(ctx, getProtocolCautionByID, id)
	var i ProtocolCaution
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
	)
	return i, err
}

const getProtocolCautionsByProtocol = `-- name: GetProtocolCautionsByProtocol :many
SELECT c.id, c.created_at, c.updated_at, c.description FROM protocol_cautions c JOIN protocol_cautions_values v ON c.id = v.caution_id WHERE v.protocol_id = $1
`

func (q *Queries) GetProtocolCautionsByProtocol(ctx context.Context, protocolID uuid.UUID) ([]ProtocolCaution, error) {
	rows, err := q.db.QueryContext(ctx, getProtocolCautionsByProtocol, protocolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProtocolCaution{}
	for rows.Next() {
		var i ProtocolCaution
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Description,
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

const getProtocolPrecautionByID = `-- name: GetProtocolPrecautionByID :one
SELECT id, created_at, updated_at, title, description FROM protocol_precautions WHERE id = $1
`

func (q *Queries) GetProtocolPrecautionByID(ctx context.Context, id uuid.UUID) (ProtocolPrecaution, error) {
	row := q.db.QueryRowContext(ctx, getProtocolPrecautionByID, id)
	var i ProtocolPrecaution
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
	)
	return i, err
}

const getProtocolPrecautionByTitleAndDescription = `-- name: GetProtocolPrecautionByTitleAndDescription :one
SELECT id, created_at, updated_at, title, description FROM protocol_precautions WHERE title = $1 AND description = $2
`

type GetProtocolPrecautionByTitleAndDescriptionParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (q *Queries) GetProtocolPrecautionByTitleAndDescription(ctx context.Context, arg GetProtocolPrecautionByTitleAndDescriptionParams) (ProtocolPrecaution, error) {
	row := q.db.QueryRowContext(ctx, getProtocolPrecautionByTitleAndDescription, arg.Title, arg.Description)
	var i ProtocolPrecaution
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
	)
	return i, err
}

const getProtocolPrecautionsByProtocol = `-- name: GetProtocolPrecautionsByProtocol :many
SELECT p.id, p.created_at, p.updated_at, p.title, p.description FROM protocol_precautions p JOIN protocol_precautions_values v ON p.id = v.precaution_id WHERE v.protocol_id = $1
`

func (q *Queries) GetProtocolPrecautionsByProtocol(ctx context.Context, protocolID uuid.UUID) ([]ProtocolPrecaution, error) {
	rows, err := q.db.QueryContext(ctx, getProtocolPrecautionsByProtocol, protocolID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ProtocolPrecaution{}
	for rows.Next() {
		var i ProtocolPrecaution
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Title,
			&i.Description,
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

const removeProtocolCautionFromProtocol = `-- name: RemoveProtocolCautionFromProtocol :exec
DELETE FROM protocol_cautions_values WHERE protocol_id = $1 AND caution_id = $2
`

type RemoveProtocolCautionFromProtocolParams struct {
	ProtocolID uuid.UUID `json:"protocol_id"`
	CautionID  uuid.UUID `json:"caution_id"`
}

func (q *Queries) RemoveProtocolCautionFromProtocol(ctx context.Context, arg RemoveProtocolCautionFromProtocolParams) error {
	_, err := q.db.ExecContext(ctx, removeProtocolCautionFromProtocol, arg.ProtocolID, arg.CautionID)
	return err
}

const removeProtocolPrecautionFromProtocol = `-- name: RemoveProtocolPrecautionFromProtocol :exec
DELETE FROM protocol_precautions_values WHERE protocol_id = $1 AND precaution_id = $2
`

type RemoveProtocolPrecautionFromProtocolParams struct {
	ProtocolID   uuid.UUID `json:"protocol_id"`
	PrecautionID uuid.UUID `json:"precaution_id"`
}

func (q *Queries) RemoveProtocolPrecautionFromProtocol(ctx context.Context, arg RemoveProtocolPrecautionFromProtocolParams) error {
	_, err := q.db.ExecContext(ctx, removeProtocolPrecautionFromProtocol, arg.ProtocolID, arg.PrecautionID)
	return err
}

const updateCautionProtocols = `-- name: UpdateCautionProtocols :exec
WITH current_protocols AS (
    SELECT pcv.protocol_id 
    FROM protocol_cautions_values pcv 
    WHERE pcv.caution_id = $1
),
to_remove AS (
    DELETE FROM protocol_cautions_values pcv
    WHERE pcv.caution_id = $1
    AND pcv.protocol_id NOT IN (SELECT unnest($2::uuid[]))
    RETURNING pcv.protocol_id
),
to_add AS (
    INSERT INTO protocol_cautions_values (caution_id, protocol_id)
    SELECT $1, new_protocol
    FROM unnest($2::uuid[]) AS new_protocol
    WHERE new_protocol NOT IN (SELECT cp.protocol_id FROM current_protocols cp)
    RETURNING protocol_id
)
SELECT 
    (SELECT COUNT(*) FROM to_remove) AS removed, 
    (SELECT COUNT(*) FROM to_add) AS added
`

type UpdateCautionProtocolsParams struct {
	CautionID   uuid.UUID   `json:"caution_id"`
	ProtocolIds []uuid.UUID `json:"protocol_ids"`
}

func (q *Queries) UpdateCautionProtocols(ctx context.Context, arg UpdateCautionProtocolsParams) error {
	_, err := q.db.ExecContext(ctx, updateCautionProtocols, arg.CautionID, pq.Array(arg.ProtocolIds))
	return err
}

const updatePrecautionProtocols = `-- name: UpdatePrecautionProtocols :exec
WITH current_protocols AS (
    SELECT pcv.protocol_id 
    FROM protocol_precautions_values pcv 
    WHERE pcv.precaution_id = $1
),
to_remove AS (
    DELETE FROM protocol_precautions_values pcv
    WHERE pcv.precaution_id = $1
    AND pcv.protocol_id NOT IN (SELECT unnest($2::uuid[]))
    RETURNING pcv.protocol_id
),
to_add AS (
    INSERT INTO protocol_precautions_values (precaution_id, protocol_id)
    SELECT $1, new_protocol
    FROM unnest($2::uuid[]) AS new_protocol
    WHERE new_protocol NOT IN (SELECT cp.protocol_id FROM current_protocols cp)
    RETURNING protocol_id
)
SELECT 
    (SELECT COUNT(*) FROM to_remove) AS removed, 
    (SELECT COUNT(*) FROM to_add) AS added
`

type UpdatePrecautionProtocolsParams struct {
	PrecautionID uuid.UUID   `json:"precaution_id"`
	ProtocolIds  []uuid.UUID `json:"protocol_ids"`
}

func (q *Queries) UpdatePrecautionProtocols(ctx context.Context, arg UpdatePrecautionProtocolsParams) error {
	_, err := q.db.ExecContext(ctx, updatePrecautionProtocols, arg.PrecautionID, pq.Array(arg.ProtocolIds))
	return err
}

const updateProtocolCaution = `-- name: UpdateProtocolCaution :one
UPDATE protocol_cautions SET description = $2 WHERE id = $1 RETURNING id, created_at, updated_at, description
`

type UpdateProtocolCautionParams struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
}

func (q *Queries) UpdateProtocolCaution(ctx context.Context, arg UpdateProtocolCautionParams) (ProtocolCaution, error) {
	row := q.db.QueryRowContext(ctx, updateProtocolCaution, arg.ID, arg.Description)
	var i ProtocolCaution
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
	)
	return i, err
}

const updateProtocolPrecaution = `-- name: UpdateProtocolPrecaution :one
UPDATE protocol_precautions SET title = $2, description = $3 WHERE id = $1 RETURNING id, created_at, updated_at, title, description
`

type UpdateProtocolPrecautionParams struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
}

func (q *Queries) UpdateProtocolPrecaution(ctx context.Context, arg UpdateProtocolPrecautionParams) (ProtocolPrecaution, error) {
	row := q.db.QueryRowContext(ctx, updateProtocolPrecaution, arg.ID, arg.Title, arg.Description)
	var i ProtocolPrecaution
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
	)
	return i, err
}

const upsertCaution = `-- name: UpsertCaution :one
WITH input_values(id, description) AS (
  VALUES (
    CASE 
      WHEN $1 = '00000000-0000-0000-0000-000000000000'::uuid 
      THEN gen_random_uuid() 
      ELSE $1 
    END,
    $2
  )
)
INSERT INTO protocol_cautions (id, description)
SELECT id, description FROM input_values
ON CONFLICT (id) DO UPDATE
SET description = EXCLUDED.description,
    updated_at = NOW()
RETURNING id, created_at, updated_at, description
`

type UpsertCautionParams struct {
	ID          interface{} `json:"id"`
	Description interface{} `json:"description"`
}

func (q *Queries) UpsertCaution(ctx context.Context, arg UpsertCautionParams) (ProtocolCaution, error) {
	row := q.db.QueryRowContext(ctx, upsertCaution, arg.ID, arg.Description)
	var i ProtocolCaution
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Description,
	)
	return i, err
}

const upsertPrecaution = `-- name: UpsertPrecaution :one
WITH input_values(id,title, description) AS (
  VALUES (
    CASE 
      WHEN $1 = '00000000-0000-0000-0000-000000000000'::uuid 
      THEN gen_random_uuid() 
      ELSE $1
    END,
    $2,
    $3
  )
)
INSERT INTO protocol_precautions (id,title, description)
SELECT id,title, description FROM input_values
ON CONFLICT (id) DO UPDATE
SET description = EXCLUDED.description,
    title = EXCLUDED.title,
    updated_at = NOW()
RETURNING id, created_at, updated_at, title, description
`

type UpsertPrecautionParams struct {
	ID          interface{} `json:"id"`
	Title       interface{} `json:"title"`
	Description interface{} `json:"description"`
}

func (q *Queries) UpsertPrecaution(ctx context.Context, arg UpsertPrecautionParams) (ProtocolPrecaution, error) {
	row := q.db.QueryRowContext(ctx, upsertPrecaution, arg.ID, arg.Title, arg.Description)
	var i ProtocolPrecaution
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Title,
		&i.Description,
	)
	return i, err
}
