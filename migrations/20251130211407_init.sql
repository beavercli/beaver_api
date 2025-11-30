-- +goose Up
-- +goose StatementBegin
CREATE TABLE tags(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    name VARCHAR(255) UNIQUE
);

CREATE TABLE languages(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    name VARCHAR(255) UNIQUE
);

CREATE TABLE contributors(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    first_name VARCHAR(255),
    last_name VARCHAR(255),
    email VARCHAR(1024) UNIQUE
);

CREATE TABLE scripts(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE,

    title VARCHAR(255) UNIQUE,
    code TEXT,
    project_url VARCHAR(1024),
    language_id UUID REFERENCES languages(id) ON DELETE SET NULL
);

CREATE TABLE script_tags(
    script_id UUID REFERENCES scripts(id) ON DELETE CASCADE,
    tag_id UUID REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY(script_id, tag_id)
);

CREATE TABLE script_contributors(
    script_id UUID REFERENCES scripts(id) ON DELETE CASCADE,
    contributor_id UUID REFERENCES contributors(id) ON DELETE CASCADE,
    PRIMARY KEY(script_id, contributor_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE script_contributors;
DROP TABLE script_tags;
DROP TABLE scripts;
DROP TABLE contributors;
DROP TABLE languages;
DROP TABLE tags;
-- +goose StatementEnd
