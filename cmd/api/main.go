package main

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/internal/shared/mailer"
	"github.com/seanflannery10/core/internal/shared/server"
	"golang.org/x/exp/slog"
)

type application struct {
	dbpool    *pgxpool.Pool
	mailer    mailer.Mailer
	secretKey []byte
	config    Config
}

func main() {
	app := &application{}

	app.init()

	if err := server.Serve(app.config.Port, app.routes()); err != nil {
		slog.Error("unable to serve application", err)
		os.Exit(exitError)
	}

	app.dbpool.Close()
}
