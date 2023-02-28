package tokens

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/validator"
)

type createAuthTokenPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p *createAuthTokenPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateEmail(v, p.Email)
	data.ValidatePasswordPlaintext(v, p.Password)

	if v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}

type createPasswordResetTokenPayload struct {
	Email string `json:"email"`
}

func (p *createPasswordResetTokenPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateEmail(v, p.Email)

	if v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}

type createActivationTokenPayload struct {
	Email string `json:"email"`
}

func (p *createActivationTokenPayload) Bind(_ *http.Request) error {
	v := validator.New()

	data.ValidateEmail(v, p.Email)

	if v.HasErrors() {
		return validator.ErrValidation
	}

	return nil
}
