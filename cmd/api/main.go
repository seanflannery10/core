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
	SMTP struct {
		Host     string `env:"SMTP_HOST,default=smtp.mailtrap.io"`
		Port     int    `env:"SMTP_PORT,default=25"`
		Username string `env:"SMTP_USERNAME"`
		Password string `env:"SMTP_PASSWORD"`
		Sender   string `env:"SMTP_SENDER,default=Test <no-reply@testdomain.com>"`
	}
	DB struct {
		DSN string `env:"DB_DSN"`
	}
}

type application struct {
	config  Config
	mailer  mailer.Mailer
	queries *data.Queries
}

func main() {
	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", helpers.GetVersion())
		os.Exit(0)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout)))

	cfg := Config{}

	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cfg.DB.DSN == "" {
		log.Fatal("DB_DSN Missing")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, cfg.DB.DSN)
	if err != nil {
		log.Fatal(err) // nolint:gocritic
	}

	defer dbpool.Close()

	slog.Info("database connection pool established")

	expvar.NewString("version").Set(helpers.GetVersion())
	expvar.Publish("goroutines", expvar.Func(func() any { return runtime.NumGoroutine() }))
	expvar.Publish("timestamp", expvar.Func(func() any { return time.Now().Unix() }))
	expvar.Publish("database", expvar.Func(func() any { return dbpool.Stat() }))

	m, err := mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender)
	if err != nil {
		log.Fatal(err, nil)
	}

	app := &application{
		config:  cfg,
		mailer:  m,
		queries: data.New(dbpool),
	}

	err = server.Serve(app.config.Connection.Port, app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
