package main

import (
	"expvar"

	"github.com/seanflannery10/core/internal/shared/errs"
	"github.com/seanflannery10/core/internal/shared/helpers"
	"github.com/seanflannery10/core/internal/shared/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() *chi.Mux {
	router := chi.NewRouter()
	env := app.env

	router.NotFound(helpers.ErrFuncWrapper(errs.ErrNotFound()))
	router.MethodNotAllowed(helpers.ErrFuncWrapper(errs.ErrMethodNotAllowed()))

	router.Use(middleware.StartSpan(env))
	router.Use(middleware.Metrics)
	router.Use(middleware.RecoverPanic)

	router.Use(cors.Handler(cors.Options{AllowedOrigins: []string{"https://*"}}))

	router.Use(middleware.Authenticate(env))

	router.Get("/debug/vars", expvar.Handler().ServeHTTP)

	return router
}
