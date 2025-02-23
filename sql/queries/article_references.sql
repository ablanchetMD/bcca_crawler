-- name: CreateArticleReference :one
INSERT INTO article_references (title, authors, journal, year, doi, pmid)
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
    doi = $6,
    pmid = $7    
WHERE id = $1
RETURNING *;

-- name: UpsertArticleReference :one
INSERT INTO article_references (id, title, authors, journal, year, doi, pmid, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
ON CONFLICT (id) DO UPDATE
SET title = EXCLUDED.title,
    authors = EXCLUDED.authors,
    journal = EXCLUDED.journal,
    year = EXCLUDED.year,
    doi = EXCLUDED.doi,
    pmid = EXCLUDED.pmid,    
    updated_at = NOW()
RETURNING *;


-- name: DeleteArticleReference :exec
DELETE FROM article_references
WHERE id = $1;

-- name: GetArticleReferenceByID :one
SELECT * FROM article_references
WHERE id = $1;

-- name: GetArticleReferenceByIDWithProtocols :one
SELECT ar.*, ARRAY_AGG(ROW(arpv.protocol_id,p.code)) AS protocol_ids
FROM article_references ar
JOIN protocol_references_value arpv ON ar.id = arpv.reference_id
JOIN protocols p ON arpv.protocol_id = p.id
WHERE ar.id = $1
GROUP BY ar.id;


-- name: GetArticleReferencesByProtocol :many
SELECT article_references.*
FROM article_references
JOIN protocol_references_value ON article_references.id = protocol_references_value.reference_id
WHERE protocol_references_value.protocol_id = $1
ORDER BY article_references.year DESC;

-- name: GetArticleReferenceByData :one
SELECT * FROM article_references
WHERE title = $1 AND authors = $2 AND journal = $3 AND year = $4;

-- name: GetArticleReferences :many
SELECT * FROM article_references
ORDER BY year DESC;

-- name: GetArticleReferencesWithProtocols :many
SELECT ar.*, ARRAY_AGG(ROW(arpv.protocol_id,p.code)) AS protocol_ids
FROM article_references ar
JOIN protocol_references_value arpv ON ar.id = arpv.reference_id
JOIN protocols p ON arpv.protocol_id = p.id
GROUP BY ar.id;

-- name: AddArticleReferenceToProtocol :exec
INSERT INTO protocol_references_value (protocol_id, reference_id)
VALUES ($1, $2);

-- name: RemoveArticleReferenceFromProtocol :exec
DELETE FROM protocol_references_value
WHERE protocol_id = $1 AND reference_id = $2;
