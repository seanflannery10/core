// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2
// source: users.sql

package data

import (
	"context"
	"time"
)

const checkUser = `-- name: CheckUser :one
SELECT EXISTS(SELECT email FROM users WHERE email = $1)::bool
`

func (q *Queries) CheckUser(ctx context.Context, email string) (bool, error) {
	row := q.db.QueryRow(ctx, checkUser, email)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, email, password_hash, activated)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, name, email, password_hash, activated, version
`

type CreateUserParams struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	PasswordHash []byte `json:"password_hash"`
	Activated    bool   `json:"activated"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser,
		arg.Name,
		arg.Email,
		arg.PasswordHash,
		arg.Activated,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
	)
	return i, err
}

const getUserFromEmail = `-- name: GetUserFromEmail :one
SELECT id, created_at, name, email, password_hash, activated, version
FROM users
WHERE email = $1
`

func (q *Queries) GetUserFromEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserFromEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
	)
	return i, err
}

const getUserFromToken = `-- name: GetUserFromToken :one
SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
FROM users
         INNER JOIN tokens
                    ON users.id = tokens.user_id
WHERE tokens.hash = $1
  AND tokens.scope = $2
  AND tokens.expiry > $3
`

type GetUserFromTokenParams struct {
	Hash   []byte    `json:"hash"`
	Scope  string    `json:"scope"`
	Expiry time.Time `json:"expiry"`
}

func (q *Queries) GetUserFromToken(ctx context.Context, arg GetUserFromTokenParams) (User, error) {
	row := q.db.QueryRow(ctx, getUserFromToken, arg.Hash, arg.Scope, arg.Expiry)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
	)
	return i, err
}

const updateUser = `-- name: UpdateUser :one
UPDATE users
SET name          = CASE WHEN $1::boolean THEN $2 ELSE name END,
    email         = CASE WHEN $3::boolean THEN $4 ELSE email END,
    password_hash = CASE WHEN $5::boolean THEN $6 ELSE password_hash END,
    activated     = CASE WHEN $7::boolean THEN $8 ELSE activated END,
    version       = version + 1
WHERE id = $9
  AND version = $10
RETURNING id, created_at, name, email, password_hash, activated, version
`

type UpdateUserParams struct {
	UpdateName         bool   `json:"update_name"`
	Name               string `json:"name"`
	UpdateEmail        bool   `json:"update_email"`
	Email              string `json:"email"`
	UpdatePasswordHash bool   `json:"update_password_hash"`
	PasswordHash       []byte `json:"password_hash"`
	UpdateActivated    bool   `json:"update_activated"`
	Activated          bool   `json:"activated"`
	ID                 int64  `json:"id"`
	Version            int32  `json:"version"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, updateUser,
		arg.UpdateName,
		arg.Name,
		arg.UpdateEmail,
		arg.Email,
		arg.UpdatePasswordHash,
		arg.PasswordHash,
		arg.UpdateActivated,
		arg.Activated,
		arg.ID,
		arg.Version,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.Name,
		&i.Email,
		&i.PasswordHash,
		&i.Activated,
		&i.Version,
	)
	return i, err
}
