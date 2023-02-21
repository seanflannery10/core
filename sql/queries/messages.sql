-- name: CreateMessage :one
INSERT INTO messages (message, user_id)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateMessage :one
UPDATE messages
SET message = $1,
    version = version + 1
WHERE id = $2
RETURNING *;

-- name: DeleteMessage :exec
DELETE
FROM messages
WHERE id = $1;

-- name: GetMessage :one
SELECT id, created_at, message, user_id, version
FROM messages
WHERE id = $1;

-- name: GetUserMessages :many
SELECT id,
       created_at,
       message,
       user_id,
       version
FROM messages
WHERE user_id = $1
ORDER BY created_at
OFFSET $2 LIMIT $3;

-- name: GetUserMessageCount :one
SELECT count(1)
FROM messages
WHERE id = $1;