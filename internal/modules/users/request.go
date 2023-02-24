package users

import (
	"errors"
	"net/http"

	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/helpers"
)

var errUserExistists = errors.New("user exists")

type Request struct {
	data.User
	Password string `json:"password"`
}

func (req *Request) Bind(r *http.Request) error {
	queries := helpers.ContextGetQueries(r)

	ok, err := queries.CheckUser(r.Context(), req.User.Email)
	if err != nil {
		return err
	}

	if ok {
		return errUserExistists
	}

	hash, err := data.GetPasswordHash(req.Password)
	if err != nil {
		return err
	}

	req.User.PasswordHash = hash

	return nil
}
