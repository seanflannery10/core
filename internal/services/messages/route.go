package messages

import (
	"github.com/go-chi/chi/v5"
	"github.com/seanflannery10/core/pkg/middleware"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {
		r.Use(middleware.RequireAuthenticatedUser)

		r.Get("/", listUserMessagesHandler)
		r.Post("/", createMessageHandler)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", showMessageHandler)
			r.Patch("/", updateMessageHandler)
			r.Delete("/", deleteMessageHandler)
		})
	})

	return r
}
