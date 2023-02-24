package users

import (
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/register", createUserHandler)
	})

	return r
}
