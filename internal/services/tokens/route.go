package tokens

import (
	"github.com/go-chi/chi/v5"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Post("/authentication", createAuthTokenHandler)
		r.Put("/activation", createActivationTokenHandler)
		r.Put("/password-reset", createPasswordResetTokenHandler)
	})

	return r
}
