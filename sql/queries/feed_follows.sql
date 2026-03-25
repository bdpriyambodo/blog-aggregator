-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (
        id, created_at, updated_at, user_id, feed_id
    )
    VALUES (
        $1, $2, $3, $4, $5
    ) 
    RETURNING *
)
SELECT 
    inserted_feed_follow.*,
    feeds.name as feed_name,
    users.name as user_name
FROM inserted_feed_follow
INNER JOIN users
ON inserted_feed_follow.user_id = users.id
INNER JOIN feeds
ON inserted_feed_follow.feed_id = feeds.id;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.*,
    users.name,
    feeds.name
FROM users
INNER JOIN feed_follows
ON users.id = feed_follows.user_id
INNER JOIN feeds
ON feed_follows.feed_id = feeds.id
WHERE users.name = $1;