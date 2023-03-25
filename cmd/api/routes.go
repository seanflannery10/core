package main

import (
	"expvar"
	"net/http"
	"net/http/pprof"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (app *application) routes() *http.ServeMux {
	// han := handlers.ServiceHandler{
	//	Mailer:  mailer.Mailer{},
	//	Queries: data.Queries{},
	//	Secret:  app.env.Config.Secret,
	//}
	//
	// srv, err := oas.NewServer(han, nil)
	// if err != nil {
	//    panic(err)
	//}

	mux := http.NewServeMux()

	// mux.Handle("/", srv)

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
