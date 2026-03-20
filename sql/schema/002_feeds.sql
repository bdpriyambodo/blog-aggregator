-- +goose Up
CREATE TABLE feeds (
    id uuid primary key,
    created_at timestamp,
    updated_at timestamp,
    name text,
    url text unique,
    user_id uuid,
    -- foreign key
    constraint fk_userid
    foreign key (user_id)
    references users (id)
    on delete cascade
);

-- +goose Down
DROP TABLE feeds;