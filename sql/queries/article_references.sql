-- name: CreateArticleReference :one
INSERT INTO article_references (title, authors, journal, year, joi, pmid)
VALUES ($1, $2, $3, $4, $5, $6)    
RETURNING *;

-- name: UpdateArticleReference :one
UPDATE article_references
SET
    updated_at = NOW(),
    title = $2,
    authors = $3,
    journal = $4,
    year = $5,
    joi = $6,
    pmid = $7    
WHERE id = $1
RETURNING *;

-- name: DeleteArticleReference :exec
DELETE FROM article_references
WHERE id = $1;

-- name: GetArticleReferenceByID :one
SELECT * FROM article_references
WHERE id = $1;

-- name: GetArticleReferencesByProtocol :many
SELECT article_references.*
FROM article_references
JOIN protocol_references_value ON article_references.id = protocol_references_value.reference_id
WHERE protocol_references_value.protocol_id = $1
ORDER BY article_references.year DESC;

-- name: AddArticleReferenceToProtocol :exec
INSERT INTO protocol_references_value (protocol_id, reference_id)
VALUES ($1, $2);

-- name: GetArticleReferenceByData :one
SELECT * FROM article_references
WHERE title = $1 AND authors = $2 AND journal = $3 AND year = $4;

-- name: AddManyArticleReferenceToProtocol :exec
INSERT INTO protocol_references_value (protocol_id, reference_id)
VALUES ($1::UUID[], $2::UUID[])
ON CONFLICT DO NOTHING;

-- name: RemoveArticleReferenceFromProtocol :exec
DELETE FROM protocol_references_value
WHERE protocol_id = $1 AND reference_id = $2;
