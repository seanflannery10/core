package users

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
)

type userResponsePayload struct {
	data.User
}

func (p userResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}

type stringResponsePayload struct {
	Message string `json:"message"`
}

func (p stringResponsePayload) Render(_ http.ResponseWriter, _ *http.Request) error {
	return nil
}
