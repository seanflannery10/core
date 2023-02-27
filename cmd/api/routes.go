package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/seanflannery10/core/internal/services/messages"
	"github.com/seanflannery10/core/internal/services/tokens"
	"github.com/seanflannery10/core/internal/services/users"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.NotFound(helpers.ErrFuncWrapper(errs.ErrNotFound))
	r.MethodNotAllowed(helpers.ErrFuncWrapper(errs.ErrMethodNotAllowed))

	r.Use(middleware.Metrics)
	r.Use(middleware.RecoverPanic)

	r.Use(middleware.SetQueriesCtx(app.queries))
	r.Use(middleware.SetMailerCtx(app.mailer))
	r.Use(middleware.Authenticate)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
	}))

	r.Get("/debug/vars", expvar.Handler().ServeHTTP)

	r.Mount("/v1/messages", messages.Router())
	r.Mount("/v1/users", users.Router())
	r.Mount("/v1/tokens", tokens.Router())

	return r
}
