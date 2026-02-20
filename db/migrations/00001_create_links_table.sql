-- +goose Up
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE links;
