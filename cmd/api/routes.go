package main

import (
	"expvar"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
	"github.com/seanflannery10/core/internal/services/messages"
	"github.com/seanflannery10/core/internal/services/tokens"
	"github.com/seanflannery10/core/internal/services/users"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/middleware"
	"net/http"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		_ = render.Render(w, r, errs.ErrNotFound)
	})

	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		_ = render.Render(w, r, errs.ErrMethodNotAllowed)
	})

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

	r.Mount("/v1/messages", messages.Router())
	r.Mount("/v1/users", users.Router())
	r.Mount("/v1/tokens", tokens.Router())

	return r
}
