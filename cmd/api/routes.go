package main

import (
	"expvar"

	"github.com/go-chi/chi/v5"
	"github.com/seanflannery10/core/internal/data"
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

	r.Use(middleware.SetQueriesCtx(data.New(app.dbpool)))
	r.Use(middleware.SetMailerCtx(app.mailer))
	r.Use(middleware.Authenticate)

	// r.Use(cors.Handler(cors.Options{
	//	AllowedOrigins: []string{"https://*", "http://*"},
	// }))

	r.Get("/debug/vars", expvar.Handler().ServeHTTP)

	r.Route("/v1/messages", func(r chi.Router) {
		r.Use(middleware.RequireAuthenticatedUser)

		r.Get("/", messages.GetMessagesUserHandler)
		r.Post("/", messages.CreateMessageHandler)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", messages.GetMessageHandler)
			r.Patch("/", messages.UpdateMessageHandler)
			r.Delete("/", messages.DeleteMessageHandler)
		})
	})

	r.Route("/v1/tokens", func(r chi.Router) {
		r.Post("/authentication", tokens.CreateTokenAuthHandler)
		r.Put("/activation", tokens.CreateTokenActivationHandler)
		r.Put("/password-reset", tokens.CreateTokenPasswordResetHandler)
	})

	r.Route("/v1/users", func(r chi.Router) {

		r.Post("/register", users.CreateUserHandler)
		r.Put("/activate", users.ActivateUserHandler)
		r.Put("/update-password", users.UpdateUserPasswordHandler)
	})

	return r
}
