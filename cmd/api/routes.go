package main

import (
	"expvar"
	"github.com/go-chi/chi/v5"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/middleware"
	"github.com/seanflannery10/core/internal/services/messages"
	"github.com/seanflannery10/core/internal/services/tokens"
	"github.com/seanflannery10/core/internal/services/users"
)

func (app *application) routes() *chi.Mux {
	r := chi.NewRouter()
	env := app.env

	r.NotFound(helpers.ErrFuncWrapper(errs.ErrNotFound))
	r.MethodNotAllowed(helpers.ErrFuncWrapper(errs.ErrMethodNotAllowed))

	r.Use(middleware.StartSpan(env))
	r.Use(middleware.Metrics)
	r.Use(middleware.RecoverPanic)

	r.Use(middleware.Authenticate(env))

	r.Get("/debug/vars", expvar.Handler().ServeHTTP)

	r.Route("/v1/messages", func(r chi.Router) {
		r.Use(middleware.RequireAuthenticatedUser)

		r.Get("/", messages.GetMessagesUserHandler(env))
		r.Post("/", messages.CreateMessageHandler(env))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", messages.GetMessageHandler(env))
			r.Patch("/", messages.UpdateMessageHandler(env))
			r.Delete("/", messages.DeleteMessageHandler(env))
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
