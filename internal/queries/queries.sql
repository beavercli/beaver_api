-- Tags

-- name: ListTags :many
SELECT * FROM tags OFFSET $1 LIMIT $2;

-- name: UpsertTag :exec
INSERT INTO tags (name) VALUES($1) ON CONFLICT (name) DO NOTHING;

-- name: GetTagIDByName :one
SELECT id FROM tags WHERE name=$1;

-- name: DeleteTagsExcept :exec
DELETE FROM tags WHERE NOT (id = ANY(sqlc.narg('ids')::BIGINT[]));

-- name: CountTags :one
SELECT COUNT(*) FROM tags;

-- Languages

-- name: ListLanguages :many
SELECT * FROM languages OFFSET $1 LIMIT $2;

-- name: UpsertLanguage :exec
INSERT INTO languages (name) VALUES($1) ON CONFLICT (name) DO NOTHING;

-- name: GetLanguageIDByName :one
SELECT id FROM languages WHERE name=$1;

-- name: DeleteLanguagesExcept :exec
DELETE FROM languages WHERE NOT (id = ANY(sqlc.narg('ids')::BIGINT[]));

-- name: CountLanguages :one
SELECT COUNT(*) FROM languages;

-- name: GetLanguageBySnippetID :one
SELECT l.* FROM languages l
INNER JOIN snippets s ON s.language_id = l.id
WHERE s.id = $1;

-- Contributors

-- name: ListContributors :many
SELECT * FROM contributors OFFSET $1 LIMIT $2;

-- name: UpsertContributor :exec
INSERT INTO contributors (first_name, last_name, email) VALUES($1, $2, $3) ON CONFLICT (email) DO NOTHING;

-- name: GetContributorIDByEmail :one
SELECT id FROM contributors WHERE email=$1;

-- name: DeleteContributorsExcept :exec
DELETE FROM contributors WHERE NOT (id = ANY(sqlc.narg('ids')::BIGINT[]));

-- name: CountContributors :one
SELECT COUNT(*) FROM contributors;

-- Snippets

-- name: UpsertSnippet :one
INSERT INTO snippets (title, code, project_url, language_id, user_id, created_at)
VALUES($1, $2, $3, $4, $5, $6)
ON CONFLICT (title) DO UPDATE SET created_at = EXCLUDED.created_at
RETURNING id;

-- name: GetSnippetIDByTitle :one
SELECT id FROM snippets WHERE title=$1;

-- name: ListUsedLanguageIDs :many
SELECT DISTINCT(language_id) FROM snippets;

-- name: DeleteSnippetsBefore :exec
DELETE FROM snippets WHERE created_at < $1;

-- name: LinkSnippetTag :exec
INSERT INTO snippet_tags (snippet_id, tag_id) VALUES($1, $2) ON CONFLICT (snippet_id, tag_id) DO NOTHING;

-- name: ListLinkedTagIDs :many
SELECT DISTINCT(tag_id) FROM snippet_tags;

-- name: LinkSnippetContributor :exec
INSERT INTO snippet_contributors (snippet_id, contributor_id) VALUES($1, $2) ON CONFLICT (snippet_id, contributor_id) DO NOTHING;

-- name: ListLinkedContributorIDs :many
SELECT DISTINCT(contributor_id) FROM snippet_contributors;

-- name: GetSnippetByID :one
SELECT
    s.id,
    s.title,
    s.code,
    s.project_url,
    s.created_at,
    s.updated_at,
    l.id AS language_id,
    l.name AS language_name
FROM snippets s
LEFT JOIN languages l ON s.language_id = l.id
WHERE s.id = $1;

-- name: GetTagsBySnippetID :many
SELECT t.id, t.name
FROM tags t
INNER JOIN snippet_tags st ON t.id = st.tag_id
WHERE st.snippet_id = $1;

-- name: GetContributorsBySnippetID :many
SELECT c.id, c.first_name, c.last_name, c.email
FROM contributors c
INNER JOIN snippet_contributors sc ON c.id = sc.contributor_id
WHERE sc.snippet_id = $1;

-- name: ListSnippetIDs :many
SELECT id FROM snippets;

-- name: ListSnippetsFiltered :many
SELECT
    s.id,
    s.title,
    s.project_url,
    l.id AS language_id,
    l.name AS language_name
FROM snippets s
LEFT JOIN languages l ON s.language_id = l.id
WHERE (sqlc.narg('language_id')::BIGINT IS NULL OR s.language_id = sqlc.narg('language_id')::BIGINT)
  AND NOT EXISTS (
    SELECT 1
    FROM (SELECT unnest(sqlc.narg('tag_ids')::BIGINT[]) AS tag_id) ft
    WHERE NOT EXISTS (
      SELECT 1
      FROM snippet_tags st
      WHERE st.snippet_id = s.id
        AND st.tag_id = ft.tag_id
    )
  )
ORDER BY s.id
OFFSET sqlc.arg('sql_offset')::INT LIMIT sqlc.arg('sql_limit')::INT;

-- name: CountSnippetsFiltered :one
SELECT COUNT(*) FROM snippets s
WHERE (sqlc.narg('language_id')::BIGINT IS NULL OR s.language_id = sqlc.narg('language_id')::BIGINT)
  AND NOT EXISTS (
    SELECT 1
    FROM (SELECT unnest(sqlc.narg('tag_ids')::BIGINT[]) AS tag_id) ft
    WHERE NOT EXISTS (
      SELECT 1
      FROM snippet_tags st
      WHERE st.snippet_id = s.id
        AND st.tag_id = ft.tag_id
    )
  );

-- name: GetTagsBySnippetIDs :many
SELECT st.snippet_id, t.id, t.name
FROM tags t
INNER JOIN snippet_tags st ON t.id = st.tag_id
WHERE st.snippet_id = ANY(sqlc.arg('snippet_ids')::BIGINT[]);

-- Users

-- name: UpsertUser :exec
INSERT INTO users (username, email, password_hash)
VALUES ($1, $2, $3)
ON CONFLICT (email) DO NOTHING;

-- name: GetUserIDByEmail :one
SELECT id FROM users WHERE email = $1;

-- name: ListAllTags :many
SELECT * FROM tags;

-- name: ListAllLanguages :many
SELECT * FROM languages;

-- name: ListAllContributors :many
SELECT * FROM contributors;

-- name: ListAllUsers :many
SELECT * FROM users;
