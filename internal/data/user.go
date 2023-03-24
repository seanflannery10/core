package data

import (
	"errors"
	"fmt"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/crypto/bcrypt"
)

const (
	PasswordCost = 13
)

var AnonymousUser = &User{} //nolint: gochecknoglobals

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
