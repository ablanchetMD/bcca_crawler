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
WITH input_values(id, title, authors, journal, year, doi, pmid) AS (
  VALUES (
    CASE 
      WHEN @id = '00000000-0000-0000-0000-000000000000'::uuid 
      THEN gen_random_uuid() 
      ELSE @id
    END,
    @title,
    @authors,
    @journal,
    @year,
    @doi,
    @pmid
  )
)
INSERT INTO article_references (id, title, authors, journal, year, doi, pmid)
SELECT id, title, authors, journal, year, doi, pmid FROM input_values
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
SELECT 
    art.*, 
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                    'id', pecv.protocol_id, 
                    'code', p.code
                )
            )
            FROM protocol_references_value pecv
            JOIN protocols p ON pecv.protocol_id = p.id
            WHERE pecv.reference_id = art.id
        ),
        '[]'
    ) AS protocol_ids
FROM 
    article_references art
WHERE art.id = $1;

-- name: DebugArticleReference :many
SELECT * FROM protocol_references_value;

-- name: GetALLArticles :many
SELECT ar.*, pr.protocol_id
FROM article_references ar
LEFT JOIN protocol_references_value pr ON ar.id = pr.reference_id
ORDER BY ar.year DESC;

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
SELECT 
    art.*, 
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                    'id', pecv.protocol_id, 
                    'code', p.code
                )
            )
            FROM protocol_references_value pecv
            JOIN protocols p ON pecv.protocol_id = p.id
            WHERE pecv.reference_id = art.id
        ),
        '[]'
    ) AS protocol_ids
FROM 
    article_references art;

-- name: GetArticleReferencesWithProtocols :many
SELECT 
    art.*, 
    COALESCE(
        (
            SELECT json_agg(
                json_build_object(
                    'id', pecv.protocol_id, 
                    'code', p.code
                )
            )
            FROM protocol_references_value pecv
            JOIN protocols p ON pecv.protocol_id = p.id
            WHERE pecv.reference_id = art.id
        ),
        '[]'
    ) AS protocol_ids
FROM 
    article_references art;

-- name: AddArticleReferenceToProtocol :exec
INSERT INTO protocol_references_value (protocol_id, reference_id)
VALUES ($1, $2);

-- name: RemoveArticleReferenceFromProtocol :exec
DELETE FROM protocol_references_value
WHERE protocol_id = $1 AND reference_id = $2;
