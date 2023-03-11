package main

import (
	"context"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-chi/docgen"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/pkg/mailer"
	"github.com/seanflannery10/core/internal/pkg/telemetry"
	"github.com/seanflannery10/core/internal/services"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/exp/slog"
)

type Config struct {
	Port         int    `env:"PORT,default=4000"`
	Env          string `env:"ENV,default=dev"`
	OTelEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT,default=api.honeycomb.io:443"`
	DSN          string `env:"DATABASE_URL,default=postgres://postgres:test@localhost:5432/test?sslmode=disable"`
	SMTP         mailer.SMTP
}

func (app *application) init() {
	generateRoutes := flag.Bool("routes", false, "Generate router documentation")
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", helpers.GetVersion())
		os.Exit(0)
	}

	if *generateRoutes {
		routesMD := []byte(docgen.MarkdownRoutesDoc(app.routes(), docgen.MarkdownOpts{
			ProjectPath: "github.com/seanflannery10/core",
			Intro:       "Routes for core API",
		}))

		err := os.WriteFile("routes.md", routesMD, 0o644) //nolint:gosec
		if err != nil {
			slog.Error("unable to create routes.md", err)
			os.Exit(1)
		}

		os.Exit(0)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)))

	err := envconfig.Process(context.Background(), &app.config)
	if err != nil {
		slog.Error("unable to process env config", err)
		os.Exit(1)
	}

	cfg := app.config

	m, err := mailer.New(cfg.SMTP)
	if err != nil {
		slog.Error("unable to create mailer", err)
		os.Exit(1)
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		slog.Error("unable to create connection pool", err)
		os.Exit(1)
	}

	tp, err := telemetry.New(cfg.OTelEndpoint, cfg.Env)
	if err != nil {
		slog.Error("unable to start telemetry", err)
		os.Exit(1)
	}

	app.dbpool = dbpool
	app.tp = tp

	app.env = services.Env{
		Queries: data.New(dbpool),
		Mailer:  m,
		Tracer:  tp.Tracer("main"),
	}

	expvar.NewString("version").Set(helpers.GetVersion())
	expvar.Publish("goroutines", expvar.Func(func() any { return runtime.NumGoroutine() }))
	expvar.Publish("timestamp", expvar.Func(func() any { return time.Now().Unix() }))
	expvar.Publish("database", expvar.Func(func() any { return dbpool.Stat() }))
}
