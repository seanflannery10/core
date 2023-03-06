package main

import (
	"expvar"

	"github.com/go-chi/chi/v5"
	"github.com/seanflannery10/core/internal/services/messages"
	"github.com/seanflannery10/core/internal/services/tokens"
	"github.com/seanflannery10/core/internal/services/users"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/middleware"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(helpers.ErrFuncWrapper(errs.ErrNotFound))
	r.MethodNotAllowed(helpers.ErrFuncWrapper(errs.ErrMethodNotAllowed))

	r.Use(middleware.Metrics)
	r.Use(middleware.RecoverPanic)

	r.Use(middleware.Authenticate(app.env))

	r.Get("/debug/vars", expvar.Handler().ServeHTTP)

	r.Route("/v1/messages", func(r chi.Router) {
		r.Use(middleware.RequireAuthenticatedUser)

		r.Get("/", messages.GetMessagesUserHandler(app.env))
		r.Post("/", messages.CreateMessageHandler(app.env))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", messages.GetMessageHandler(app.env))
			r.Patch("/", messages.UpdateMessageHandler(app.env))
			r.Delete("/", messages.DeleteMessageHandler(app.env))
		})
	})

	r.Route("/v1/tokens", func(r chi.Router) {
		r.Post("/authentication", tokens.CreateTokenAuthHandler(app.env))
		r.Put("/activation", tokens.CreateTokenActivationHandler(app.env))
		r.Put("/password-reset", tokens.CreateTokenPasswordResetHandler(app.env))
	})

	r.Route("/v1/users", func(r chi.Router) {
		r.Post("/register", users.CreateUserHandler(app.env))
		r.Put("/activate", users.ActivateUserHandler(app.env))
		r.Put("/update-password", users.UpdateUserPasswordHandler(app.env))
	})

	return r
}
