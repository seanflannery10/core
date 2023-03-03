package main

import (
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/go-chi/docgen"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/mailer"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/exp/slog"
)

type Config struct {
	Connection struct {
		Port int    `env:"PORT,default=4000"`
		Env  string `env:"ENV,default=dev"`
	}
	DB struct {
		DSN string `env:"DB_DSN,default=postgres://postgres:test@localhost:5432/test?sslmode=disable"`
	}
	SMTP mailer.SMTP
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

	err := envconfig.Process(ctx, &app.config)
	if err != nil {
		slog.Error("unable to process env config", err)
		os.Exit(1)
	}

	m, err := mailer.New(app.config.SMTP)
	if err != nil {
		slog.Error("unable to create mailer", err)
		os.Exit(1)
	}

	dbpool, err := pgxpool.New(ctx, app.config.DB.DSN)
	if err != nil {
		slog.Error("unable to create connection pool", err)
		os.Exit(1)
	}

	app.dbpool = dbpool
	app.mailer = m

	expvar.NewString("version").Set(helpers.GetVersion())
	expvar.Publish("goroutines", expvar.Func(func() any { return runtime.NumGoroutine() }))
	expvar.Publish("timestamp", expvar.Func(func() any { return time.Now().Unix() }))
	expvar.Publish("database", expvar.Func(func() any { return dbpool.Stat() }))
}
