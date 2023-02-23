package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/seanflannery10/core/internal/httperrors"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.NotFound(httperrors.NotFound)
	r.MethodNotAllowed(httperrors.MethodNotAllowed)

	r.Use(app.metrics)
	r.Use(app.recoverPanic)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))

	r.Use(app.authenticate)

	r.Get("/debug/vars", expvar.Handler().ServeHTTP)
	r.Get("/healthcheck", app.healthCheckHandler)

	r.Route("/v1/messages", func(r chi.Router) {
		r.Use(app.requireAuthenticatedUser)

		r.Get("/", app.listUserMessagesHandler)
		r.Post("/", app.createMessageHandler)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.showMessageHandler)
			r.Patch("/", app.updateMessageHandler)
			r.Delete("/", app.deleteMessageHandler)
		})
	})

	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/register", app.registerUserHandler)
		r.Put("/activate", app.activateUserHandler)
		r.Put("/update-password", app.updateUserPasswordHandler)
	})

	r.Route("/v1/tokens", func(r chi.Router) {
		r.Post("/authentication", app.createAuthenticationTokenHandler)
		r.Put("/activation", app.createActivationTokenHandler)
		r.Put("/password-reset", app.createPasswordResetTokenHandler)
	})

	return r
}