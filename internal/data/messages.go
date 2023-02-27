package data

import (
	"net/http"

	"github.com/seanflannery10/core/pkg/validator"
)

func (m *Message) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

func ValidateMessage(v *validator.Validator, message string) {
	v.Check(message != "", "message", "must be provided")
	v.Check(len(message) <= 512, "message", "must not be more than 512 bytes long")
}
