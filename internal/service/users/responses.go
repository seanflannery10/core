package users

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
)

type userResponse struct {
	data.User
}

func (res userResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

type messageResponse struct {
	Message string `json:"message"`
}

func (res messageResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
