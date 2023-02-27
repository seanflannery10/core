package users

import (
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/register", createUserHandler)
		r.Put("/activate", activateUserHandler)
		r.Put("/update-password", updateUserPasswordHandler)
	})

	return r
}
