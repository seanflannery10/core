package data

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/seanflannery10/core/internal/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordCost  = 13
	keyEmail      = "email"
	keyPassword   = "password"
	keyName       = "name"
	requiredField = "must be provided"
)

var AnonymousUser = &User{} //nolint: gochecknoglobals

func (u *User) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func (u *User) IsAnonymous() bool {
	return cmp.Equal(u, AnonymousUser)
}

func (u *User) SetPassword(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), PasswordCost)
	if err != nil {
		return fmt.Errorf("failed set password: %w", err)
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
			return false, fmt.Errorf("failed compare password: %w", err)
		}
	}

	return true, nil
}

func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", keyEmail, requiredField)
	v.Check(validator.Matches(email, validator.RgxEmail), keyEmail, "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
	v.Check(password != "", keyPassword, requiredField)
	v.Check(len(password) >= 8, keyPassword, "must be at least 8 bytes long")
	v.Check(len(password) <= 72, keyPassword, "must not be more than 72 bytes long")
}

func ValidateName(v *validator.Validator, name string) {
	v.Check(name != "", keyName, requiredField)
	v.Check(len(name) <= 500, keyName, "must not be more than 500 bytes long")
}

func (q *Queries) GetUserFromTokenHelper(ctx context.Context, tokenPlaintext, scope string) (User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	user, err := q.GetUserFromToken(ctx, GetUserFromTokenParams{
		Hash:   tokenHash[:],
		Scope:  scope,
		Expiry: time.Now(),
	})
	if err != nil {
		return User{}, fmt.Errorf("failed get user from token: %w", err)
	}

	return user, nil
}
