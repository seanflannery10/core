package users

import (
	"net/http"

	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type createUserPayload struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash []byte
}

func (p *createUserPayload) Bind(r *http.Request) error {
	v := validator.New()

	validateName(v, p.Name)
	validateEmail(v, p.Email)
	validatePasswordPlaintext(v, p.Password)

	if v.HasErrors() {
		return validator.ValidationError{Errors: v.Errors}
	}

	queries := helpers.ContextGetQueries(r)

	ok, err := queries.CheckUser(r.Context(), p.Email)
	if err != nil {
		return err
	}

	if ok {
		v.AddError("email", "a user with this email address already exists")
		return validator.ValidationError{Errors: v.Errors}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(p.Password), 14)
	if err != nil {
		return err
	}

	p.PasswordHash = hash

	return nil
}

type activateUserPayload struct {
	TokenPlaintext string `json:"token"`
}

func (p *activateUserPayload) Bind(_ *http.Request) error {
	v := validator.New()

	validateTokenPlaintext(v, p.TokenPlaintext)

	if v.HasErrors() {
		return validator.ValidationError{Errors: v.Errors}
	}

	return nil
}

type updateUserPasswordPayload struct {
	Password       string `json:"password"`
	TokenPlaintext string `json:"token"`
}

func (p *updateUserPasswordPayload) Bind(_ *http.Request) error {
	v := validator.New()

	validatePasswordPlaintext(v, p.Password)
	validateTokenPlaintext(v, p.TokenPlaintext)

	if v.HasErrors() {
		return validator.ValidationError{Errors: v.Errors}
	}

	return nil
}
