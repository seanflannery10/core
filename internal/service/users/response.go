package users

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
)

type UserResponse struct {
	data.User
}

func NewUserResponse(user data.User) UserResponse {
	return UserResponse{User: user}
}

func (res UserResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
