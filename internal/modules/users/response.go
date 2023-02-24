package users

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
)

type Response struct {
	data.User
}

func NewResponse(user data.User) Response {
	return Response{User: user}
}

func (res Response) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
