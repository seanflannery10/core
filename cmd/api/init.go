package main

import (
	"context"
	"encoding/hex"
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

const (
	exitGood  = 0
	exitError = 1
	filePerm  = 644
)

func (app *application) init() {
	generateRoutes := flag.Bool("routes", false, "Generate router documentation")
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		_, _ = fmt.Printf("Version:\t%s\n", helpers.GetVersion()) //nolint:forbidigo
		os.Exit(exitGood)
	}

	if *generateRoutes {
		routesMD := []byte(docgen.MarkdownRoutesDoc(app.routes(), docgen.MarkdownOpts{
			ProjectPath: "github.com/seanflannery10/core",
			Intro:       "Routes for core API",
		}))

		err := os.WriteFile("routes.md", routesMD, filePerm)
		if err != nil {
			slog.Error("unable to create routes.md", err)
			os.Exit(exitError)
		}

		os.Exit(exitGood)
	}

	secret, err := hex.DecodeString(os.Getenv("SECRET_KEY"))
	if err != nil || os.Getenv("SECRET_KEY") == "" {
		slog.Error("unable to decode sercret", err)
		os.Exit(exitError)
	}

	app.config.Secret = secret

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)))

	err = envconfig.Process(context.Background(), &app.config)
	if err != nil {
		slog.Error("unable to process env config", err)
		os.Exit(exitError)
	}

	cfg := app.config

	mail, err := mailer.New(cfg.SMTP)
	if err != nil {
		slog.Error("unable to create mailer", err)
		os.Exit(exitError)
	}

	dbpool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("unable to create connection pool", err)
		os.Exit(exitError)
	}

	tracerProvider, err := telemetry.New(cfg.OTelEndpoint, cfg.Env)
	if err != nil {
		slog.Error("unable to start telemetry", err)
		os.Exit(exitError)
	}

	app.dbpool = dbpool
	app.tp = tracerProvider

	app.env = services.Env{
		Queries: data.New(dbpool),
		Mailer:  mail,
		Tracer:  tracerProvider.Tracer("main"),
	}

	expvar.NewString("version").Set(helpers.GetVersion())
	expvar.Publish("goroutines", expvar.Func(func() any { return runtime.NumGoroutine() }))
	expvar.Publish("timestamp", expvar.Func(func() any { return time.Now().Unix() }))
	expvar.Publish("database", expvar.Func(func() any { return dbpool.Stat() }))
}
