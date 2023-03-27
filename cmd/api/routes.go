package main

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/handler"
	"github.com/seanflannery10/core/internal/oas"
	"github.com/seanflannery10/core/internal/shared/mailer"
	"github.com/seanflannery10/core/internal/shared/middleware"
)

func (app *application) routes() *http.ServeMux {
	h := &handler.Handler{
		Mailer:  mailer.Mailer{},
		Queries: data.Queries{},
		Secret:  app.config.Secret,
	}

	srv, err := oas.NewServer(h, oas.WithMiddleware(middleware.Authenticate(h.Queries)), oas.WithErrorHandler(middleware.ErrorHandler))
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.Handle("/", srv)

	mux.Handle("/metrics", promhttp.Handler())

	// Register pprof handlers.
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	mux.HandleFunc("/debug/vars", expvar.Handler().ServeHTTP)

	return mux
}
