-- name: CreateToken :one
INSERT INTO tokens (hash, user_id, expiry, scope)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CheckToken :one
SELECT EXISTS(SELECT 1
              FROM tokens
              WHERE active = false
                AND scope = $1
                AND hash = $2
                AND user_id = $3)::bool;

-- name: DeactivateToken :exec
UPDATE tokens
SET active = false
WHERE scope = $1
  AND hash = $2
  AND user_id = $3;

-- name: DeleteTokens :exec
DELETE
FROM tokens
WHERE scope = $1
  AND user_id = $2;
