-- +goose Up
-- +goose StatementBegin
CREATE TABLE tags(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    name VARCHAR(255) UNIQUE
);

CREATE TABLE languages(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    name VARCHAR(255) UNIQUE
);

CREATE TABLE contributors(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(1024) UNIQUE
);

CREATE TABLE snippets(
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    title VARCHAR(255) UNIQUE,
    code TEXT,
    project_url VARCHAR(1024),
    language_id BIGINT REFERENCES languages(id) ON DELETE SET NULL
);

CREATE TABLE snippet_tags(
    snippet_id BIGINT REFERENCES snippets(id) ON DELETE CASCADE,
    tag_id BIGINT REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY(snippet_id, tag_id)
);

CREATE TABLE snippet_contributors(
    snippet_id BIGINT REFERENCES snippets(id) ON DELETE CASCADE,
    contributor_id BIGINT REFERENCES contributors(id) ON DELETE CASCADE,
    PRIMARY KEY(snippet_id, contributor_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE snippet_contributors;
DROP TABLE snippet_tags;
DROP TABLE snippets;
DROP TABLE contributors;
DROP TABLE languages;
DROP TABLE tags;
-- +goose StatementEnd
