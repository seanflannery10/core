package main

import (
	"expvar"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/seanflannery10/core/internal/pkg/errs"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/middleware"
	"github.com/seanflannery10/core/internal/services/messages"
	"github.com/seanflannery10/core/internal/services/tokens"
	"github.com/seanflannery10/core/internal/services/users"
)

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()
	env := &app.env

	router.NotFound(helpers.ErrFuncWrapper(errs.ErrNotFound()))
	router.MethodNotAllowed(helpers.ErrFuncWrapper(errs.ErrMethodNotAllowed()))

	router.Use(middleware.StartSpan(env))
	router.Use(middleware.Metrics)
	router.Use(middleware.RecoverPanic)

	router.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"https://*"}}))

	router.Use(middleware.Authenticate(env))

	router.Get("/debug/vars", expvar.Handler().ServeHTTP)

	router.Route("/v1/messages", func(r chi.Router) {
		r.Use(middleware.RequireAuthenticatedUser)

		r.Get("/", messages.GetMessagesUserHandler(env))
		r.Post("/", messages.CreateMessageHandler(env))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", messages.GetMessageHandler(env))
			r.Put("/", messages.UpdateMessageHandler(env))
			r.Delete("/", messages.DeleteMessageHandler(env))
		})
	})

	router.Route("/v1/tokens", func(r chi.Router) {
		r.Post("/access", tokens.CreateTokenAccessHandler(env))
		r.Post("/activation", tokens.CreateTokenActivationHandler(env))
		r.Post("/password-reset", tokens.CreateTokenPasswordResetHandler(env))
		r.Post("/refresh", tokens.CreateTokenRefreshHandler(env))
	})

	router.Route("/v1/users", func(r chi.Router) {
		r.Post("/register", users.CreateUserHandler(env))
		r.Patch("/activate", users.ActivateUserHandler(env))
		r.Patch("/update-password", users.UpdateUserPasswordHandler(env))
	})

	return router
}
