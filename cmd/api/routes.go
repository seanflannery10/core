package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/seanflannery10/core/internal/modules/users"
	"github.com/seanflannery10/core/pkg/httperrors"
	"github.com/seanflannery10/core/pkg/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.NotFound(httperrors.NotFound)
	r.MethodNotAllowed(httperrors.MethodNotAllowed)

	r.Use(middleware.Metrics)
	r.Use(middleware.RecoverPanic)

	r.Use(middleware.SetQueriesCtx(app.queries))
	r.Use(middleware.SetMailerCtx(app.mailer))
	r.Use(middleware.Authenticate)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))

	r.Get("/debug/vars", expvar.Handler().ServeHTTP)
	r.Get("/healthcheck", app.healthCheckHandler)

	r.Route("/v1/messages", func(r chi.Router) {
		r.Use(middleware.RequireAuthenticatedUser)

		r.Get("/", app.listUserMessagesHandler)
		r.Post("/", app.createMessageHandler)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", app.showMessageHandler)
			r.Patch("/", app.updateMessageHandler)
			r.Delete("/", app.deleteMessageHandler)
		})
	})

	r.Mount("/", users.Router())

	r.Route("/v1/tokens", func(r chi.Router) {
		r.Post("/authentication", app.createAuthenticationTokenHandler)
		r.Put("/activation", app.createActivationTokenHandler)
		r.Put("/password-reset", app.createPasswordResetTokenHandler)
	})

	return r
}
