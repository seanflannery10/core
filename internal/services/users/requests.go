package users

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

type createUserPayload struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash []byte
	v            *validator.Validator
}

func (p *createUserPayload) Bind(r *http.Request) error {
	data.ValidateName(p.v, p.Name)
	data.ValidateEmail(p.v, p.Email)
	data.ValidatePasswordPlaintext(p.v, p.Password)

	if p.v.HasErrors() {
		return validator.ErrValidation
	}

	queries := helpers.ContextGetQueries(r)

	ok, err := queries.CheckUser(r.Context(), p.Email)
	if err != nil {
		return err
	}

	if ok {
		p.v.AddError("email", "a user with this email address already exists")
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
	v              *validator.Validator
}

func (p *activateUserPayload) Bind(_ *http.Request) error {
	data.ValidateTokenPlaintext(p.v, p.TokenPlaintext)

	if p.v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}

type updateUserPasswordPayload struct {
	Password       string `json:"password"`
	TokenPlaintext string `json:"token"`
	v              *validator.Validator
}

func (p *updateUserPasswordPayload) Bind(_ *http.Request) error {
	data.ValidatePasswordPlaintext(p.v, p.Password)
	data.ValidateTokenPlaintext(p.v, p.TokenPlaintext)

	if p.v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}
