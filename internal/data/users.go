package data

import (
	"context"
	"crypto/sha256"
	"errors"
	"net/http"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seanflannery10/core/internal/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

func (u *User) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (u *User) IsAnonymous() bool {
	return cmp.Equal(u, AnonymousUser)
}

func (u *User) SetPassword(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 14)
	if err != nil {
		return err
	}

	u.PasswordHash = hash

	return nil
}

func (u *User) ComparePasswords(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.RgxEmail), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", "password", "must be provided")
	v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
	v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}

func ValidateName(v *validator.Validator, name string) {
	v.Check(name != "", "name", "must be provided")
	v.Check(len(name) <= 500, "name", "must not be more than 500 bytes long")
}

func (q *Queries) GetUserFromTokenHelper(ctx context.Context, tokenPlaintext, scope string) (User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	user, err := q.GetUserFromToken(ctx, GetUserFromTokenParams{
		Hash:   tokenHash[:],
		Scope:  scope,
		Expiry: pgtype.Timestamptz{Time: time.Now(), Valid: true},
	})
	if err != nil {
		return User{}, err
	}

	return user, nil
}
