-- +goose Up
CREATE TABLE feed_follows (
    id uuid primary key,
    created_at timestamp,
    updated_at timestamp,
    user_id uuid,
    feed_id uuid,
    -- foreign key
    constraint fk_userid
    foreign key (user_id)
    references users (id)
    on delete cascade,
    -- foreign key
    constraint fk_feedid
    foreign key (feed_id)
    references feeds (id)
    on delete cascade,

    CONSTRAINT unique_user_feed UNIQUE (user_id, feed_id)
);


-- +goose Down
DROP TABLE feed_follows;