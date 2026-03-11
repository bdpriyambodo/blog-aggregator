-- +goose Up
CREATE TABLE users (
    id uuid primary key,
    created_at timestamp,
    updated_at timestamp,
    name text not null unique
);

-- +goose Down
DROP TABLE users;