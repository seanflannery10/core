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
	"github.com/seanflannery10/core/internal/shared/mailer"
	"github.com/seanflannery10/core/internal/shared/utils"
	"github.com/sethvargo/go-envconfig"
	"golang.org/x/exp/slog"
)

const (
	exitGood  = 0
	exitError = 1
)

type Config struct {
	Env          string `env:"ENV,default=dev"`
	OTelEndpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT,default=api.honeycomb.io:443"`
	DatabaseURL  string `env:"DATABASE_URL,default=postgres://postgres:test@localhost:5432/test?sslmode=disable"`
	SMTP         mailer.SMTP
	Port         int32 `env:"PORT,default=4000"`
}

func (app *application) init() {
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		_, _ = fmt.Printf("Version:\t%s\n", utils.GetVersion()) //nolint:forbidigo
		os.Exit(exitGood)
	}

	cfg := &Config{}

	err := envconfig.Process(context.Background(), cfg)
	if err != nil {
		slog.Error("unable to process env config", err)
		os.Exit(exitError)
	}

	secretString := os.Getenv("SECRET_KEY")
	if secretString == "" {
		slog.Error("secret key not set", err)
		os.Exit(exitError)
	}

	secret, err := hex.DecodeString(secretString)
	if err != nil {
		slog.Error("unable to decode secret key", err)
		os.Exit(exitError)
	}

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

	app.secretKey = secret
	app.dbpool = dbpool
	app.mailer = mail
	app.config = *cfg

	expvar.NewString("version").Set(utils.GetVersion())
	expvar.Publish("goroutines", expvar.Func(func() any { return runtime.NumGoroutine() }))
	expvar.Publish("timestamp", expvar.Func(func() any { return time.Now().Unix() }))
	expvar.Publish("database", expvar.Func(func() any { return dbpool.Stat() }))
}
