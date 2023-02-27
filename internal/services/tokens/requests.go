package tokens

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/validator"
)

type createAuthTokenPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	v        *validator.Validator
}

func (p *createAuthTokenPayload) Bind(_ *http.Request) error {
	data.ValidateEmail(p.v, p.Email)
	data.ValidatePasswordPlaintext(p.v, p.Password)

	if p.v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}

type createPasswordResetTokenPayload struct {
	Email string `json:"email"`
	v     *validator.Validator
}

func (p *createPasswordResetTokenPayload) Bind(_ *http.Request) error {
	data.ValidateEmail(p.v, p.Email)

	if p.v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}

type createActivationTokenPayload struct {
	Email string `json:"email"`
	v     *validator.Validator
}

func (p *createActivationTokenPayload) Bind(_ *http.Request) error {
	data.ValidateEmail(p.v, p.Email)

	if p.v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}
