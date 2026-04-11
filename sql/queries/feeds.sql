-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeed :one
SELECT * FROM feeds WHERE name = $1;

-- name: GetFeedUrl :one
SELECT * FROM feeds WHERE url = $1;

-- name: DeleteAllFeeds :exec
DELETE from feeds;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name
from feeds
left join users 
on feeds.user_id = users.id;

-- name: MarkFeedFetched :exec
UPDATE feeds 
SET last_fetched_at = NOW(), updated_at = NOW()
where id = $1;

-- name: GetNextFeedToFetch :one
SELECT id, url FROM feeds
ORDER BY last_fetched_at asc NULLS first
LIMIT 1
;
