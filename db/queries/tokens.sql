-- name: CreateToken :one
INSERT INTO tokens (hash, user_id, expiry, scope)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CreateRefreshToken :one
INSERT INTO tokens (hash, user_id, expiry, session, scope)
VALUES ($1, $2, $3, $4, "refresh")
RETURNING *;

-- name: DeactivateRefreshTokens :exec
UPDATE tokens
SET active = false
WHERE session = $1
  AND user_id = $2;

-- name: DeleteSessionTokensForUser :exec
DELETE
FROM tokens
WHERE scope = "refresh"
  AND user_id = $1
  AND session = $2;


-- name: DeleteAllTokensForUser :exec
DELETE
FROM tokens
WHERE scope = $1
  AND user_id = $2;