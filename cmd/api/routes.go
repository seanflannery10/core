package main

import (
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/seanflannery10/core/internal/generated/api"
	"github.com/seanflannery10/core/internal/generated/data"
	"github.com/seanflannery10/core/internal/server/handler"
	"golang.org/x/exp/slog"
)

func (app *application) routes() *http.ServeMux {
	newHandler := &handler.Handler{
		Queries: data.New(app.dbpool),
		Mailer:  app.mailer,
		Secret:  app.secretKey,
	}

	srv, err := api.NewServer(
		newHandler,
		&security{Queries: data.New(app.dbpool), SecretKey: app.secretKey},
		api.WithMiddleware(app.RecoverPanic()),
		api.WithErrorHandler(handler.ErrorHandler),
	)
	if err != nil {
		slog.Error("unable to create new server", err)
		os.Exit(exitError) //nolint:revive
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
