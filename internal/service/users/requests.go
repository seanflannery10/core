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
	Validator    *validator.Validator
}

func (p *createUserPayload) Bind(r *http.Request) error {
	validateName(p.Validator, p.Name)
	validateEmail(p.Validator, p.Email)
	validatePasswordPlaintext(p.Validator, p.Password)

	if p.Validator.HasErrors() {
		return validator.ErrValidation
	}

	queries := helpers.ContextGetQueries(r)

	ok, err := queries.CheckUser(r.Context(), p.Email)
	if err != nil {
		return err
	}

	if ok {
		p.Validator.AddError("email", "a user with this email address already exists")
		return validator.ErrValidation
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
	Validator      *validator.Validator
}

func (p *activateUserPayload) Bind(_ *http.Request) error {
	validateTokenPlaintext(p.Validator, p.TokenPlaintext)

	if p.Validator.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}

type updateUserPasswordPayload struct {
	Password       string `json:"password"`
	TokenPlaintext string `json:"token"`
	Validator      *validator.Validator
}

func (p *updateUserPasswordPayload) Bind(_ *http.Request) error {
	validatePasswordPlaintext(p.Validator, p.Password)
	validateTokenPlaintext(p.Validator, p.TokenPlaintext)

	if p.Validator.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}
