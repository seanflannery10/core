package main

import (
	"context"
	"expvar"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/helpers"
	"github.com/seanflannery10/core/internal/mailer"
	"github.com/sethvargo/go-envconfig"
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
		Sender   string `env:"SMTP_SENDER,default=Greenlight <no-reply@testdomain.com>"`
	}
	DB struct {
		DSN string `env:"DB_DSN"`
	}
}

type application struct {
	config  Config
	mailer  mailer.Mailer
	queries *data.Queries
	wg      sync.WaitGroup
}

func main() {
	cfg := Config{}

	err := envconfig.Process(context.Background(), &cfg)
	if err != nil {
		log.Fatal(err)
	}

	displayVersion := flag.Bool("version", false, "Display version and exit")
	flag.Parse()

	if *displayVersion {
		fmt.Printf("Version:\t%s\n", helpers.GetVersion())
		os.Exit(0)
	}

	m, err := mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender)
	if err != nil {
		log.Fatal(err, nil)
	}

	dbpool, err := helpers.NewDBPool(cfg.DB.DSN)
	if err != nil {
		log.Fatal(err, nil)
	}

	expvar.NewString("version").Set(helpers.GetVersion())

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	expvar.Publish("database", expvar.Func(func() any {
		return dbpool.Stat()
	}))

	queries := data.New(dbpool)

	app := &application{
		config:  cfg,
		mailer:  m,
		queries: queries,
	}

	err = app.serve()
	if err != nil {
		log.Fatal(err)
	}
}
