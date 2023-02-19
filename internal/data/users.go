package data

import (
	"errors"
	"fmt"

	"github.com/seanflannery10/core/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

var AnonymousUser = &User{}

func (u *User) IsAnonymous() bool {
	return u == AnonymousUser
}

func (u *User) SetPassword(plaintextPassword string) error {
	hash, err := GetPasswordHash(plaintextPassword)
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

func GetPasswordHash(plaintextPassword string) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 14)
	if err != nil {
		return nil, err
	}

	return hash, nil
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

func ValidateNewUserParams(v *validator.Validator, user CreateUserParams) {
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

	ValidateEmail(v, fmt.Sprintf("%v", user.Email))
}
