package main

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/internal/shared/mailer"
	"golang.org/x/exp/slog"
)

type application struct {
	dbpool *pgxpool.Pool
	config *Config
	mailer mailer.Mailer
}

func main() {
	app := &application{}

	app.init()

	if err := serve(app.config.Port, app.routes()); err != nil {
		slog.Error("unable to serve application", err)
		os.Exit(exitError)
	}

	app.shutdown()
}
