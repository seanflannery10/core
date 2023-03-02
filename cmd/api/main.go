package main

import (
	"context"
	"expvar"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/go-chi/docgen"
	_ "github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-launcher-go/launcher"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/mailer"
	"github.com/seanflannery10/core/pkg/server"
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

type application struct {
	Config  Config
	Mailer  mailer.Mailer
	Queries *data.Queries
}

func main() {
	app := &application{}

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
			log.Fatal(err)
		}

		os.Exit(0)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)))

	otelShutdown, err := launcher.ConfigureOpenTelemetry()
	if err != nil {
		log.Fatalf("error setting up OTel SDK - %e", err)
	}

	defer otelShutdown()

	err = envconfig.Process(context.Background(), &app.Config)
	if err != nil {
		log.Fatal(err) // nolint:gocritic
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, app.Config.DB.DSN)
	if err != nil {
		log.Fatal(err) // nolint:gocritic
	}

	defer dbpool.Close()

	slog.Info("database connection pool established")

	expvar.NewString("version").Set(helpers.GetVersion())
	expvar.Publish("goroutines", expvar.Func(func() any { return runtime.NumGoroutine() }))
	expvar.Publish("timestamp", expvar.Func(func() any { return time.Now().Unix() }))
	expvar.Publish("database", expvar.Func(func() any { return dbpool.Stat() }))

	m, err := mailer.New(app.Config.SMTP)
	if err != nil {
		log.Fatal(err, nil)
	}

	app.Mailer = m
	app.Queries = data.New(dbpool)

	err = server.Serve(app.Config.Connection.Port, app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
