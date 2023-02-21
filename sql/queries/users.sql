-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, activated)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: CheckUser :one
SELECT EXISTS(SELECT email FROM users WHERE email = $1)::bool;

-- name: UpdateUser :one
UPDATE users
SET name          = CASE WHEN @update_name::boolean THEN @name ELSE name END,
    email         = CASE WHEN @update_email::boolean THEN @email ELSE email END,
    password_hash = CASE WHEN @update_password_hash::boolean THEN @password_hash ELSE password_hash END,
    activated     = CASE WHEN @update_activated::boolean THEN @activated ELSE activated END,
    version       = version + 1
WHERE id = @id
  AND version = @version
RETURNING *;

-- name: GetUserFromEmail :one
SELECT id, created_at, name, email, password_hash, activated, version
FROM users
WHERE email = $1;

-- name: GetUserFromToken :one
SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
FROM users
         INNER JOIN tokens
                    ON users.id = tokens.user_id
WHERE tokens.hash = $1
  AND tokens.scope = $2
  AND tokens.expiry > $3;