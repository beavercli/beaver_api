-- Tags

-- name: ListTags :many
SELECT * FROM tags OFFSET $1 LIMIT $2;

-- name: UpsertTag :exec
INSERT INTO tags (name) VALUES($1) ON CONFLICT (name) DO NOTHING;

-- name: GetTagIDByName :one
SELECT id FROM tags WHERE name=$1;

-- name: DeleteTagsExcept :exec
DELETE FROM tags WHERE NOT (id = ANY(sqlc.narg('ids')::UUID[]));

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
DELETE FROM languages WHERE NOT (id = ANY(sqlc.narg('ids')::UUID[]));

-- name: CountLanguages :one
SELECT COUNT(*) FROM languages;

-- Contributors

-- name: ListContributors :many
SELECT * FROM contributors OFFSET $1 LIMIT $2;

-- name: UpsertContributor :exec
INSERT INTO contributors (first_name, last_name, email) VALUES($1, $2, $3) ON CONFLICT (email) DO NOTHING;

-- name: GetContributorIDByEmail :one
SELECT id FROM contributors WHERE email=$1;

-- name: DeleteContributorsExcept :exec
DELETE FROM contributors WHERE NOT (id = ANY(sqlc.narg('ids')::UUID[]));

-- Scripts

-- name: UpsertScript :exec
INSERT INTO scripts (title, code, project_url, language_id, created_at) VALUES($1, $2, $3, $4, $5)
ON CONFLICT (title) DO UPDATE SET created_at = EXCLUDED.created_at;

-- name: GetScriptIDByTitle :one
SELECT id FROM scripts WHERE title=$1;

-- name: ListUsedLanguageIDs :many
SELECT DISTINCT(language_id) FROM scripts;

-- name: DeleteScriptsBefore :exec
DELETE FROM scripts WHERE created_at < $1;

-- name: LinkScriptTag :exec
INSERT INTO script_tags (script_id, tag_id) VALUES($1, $2) ON CONFLICT (script_id, tag_id) DO NOTHING;

-- name: ListLinkedTagIDs :many
SELECT DISTINCT(tag_id) FROM script_tags;

-- name: LinkScriptContributor :exec
INSERT INTO script_contributors (script_id, contributor_id) VALUES($1, $2) ON CONFLICT (script_id, contributor_id) DO NOTHING;

-- name: ListLinkedContributorIDs :many
SELECT DISTINCT(contributor_id) FROM script_contributors;
