package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/pkg/mailer"
	"github.com/seanflannery10/core/pkg/server"
	"golang.org/x/exp/slog"
)

type application struct {
	config Config
	mailer mailer.Mailer
	dbpool *pgxpool.Pool
}

var ctx = context.Background()

func main() {
	app := &application{}

	app.init()
	app.otel()

	err := server.Serve(app.config.Connection.Port, app.routes())
	if err != nil {
		slog.Error("unable to serve application", err)
		os.Exit(1)
	}

	app.dbpool.Close()
}
