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
)

func (app *application) init() {
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		_, _ = fmt.Printf("Version:\t%s\n", helpers.GetVersion()) //nolint:forbidigo
		os.Exit(exitGood)
	}

	cfg := &services.Config{}

	err := envconfig.Process(context.Background(), cfg)
	if err != nil {
		slog.Error("unable to process env config", err)
		os.Exit(exitError)
	}

	secret, err := hex.DecodeString(os.Getenv("SECRET_KEY"))
	if err != nil || os.Getenv("SECRET_KEY") == "" {
		slog.Error("unable to decode sercret", err)
		os.Exit(exitError)
	}

	cfg.Secret = secret

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)))

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

	app.env = &services.Env{
		Queries: data.New(dbpool),
		Mailer:  mail,
		Tracer:  tracerProvider.Tracer("main"),
		Config:  cfg,
	}

	expvar.NewString("version").Set(helpers.GetVersion())
	expvar.Publish("goroutines", expvar.Func(func() any { return runtime.NumGoroutine() }))
	expvar.Publish("timestamp", expvar.Func(func() any { return time.Now().Unix() }))
	expvar.Publish("database", expvar.Func(func() any { return dbpool.Stat() }))
}
