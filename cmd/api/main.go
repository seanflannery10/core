package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/pkg/mailer"
	"github.com/seanflannery10/core/pkg/server"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/exp/slog"
)

type application struct {
	config Config
	mailer mailer.Mailer
	dbpool *pgxpool.Pool
	tp     *trace.TracerProvider
}

func main() {
	app := &application{}

	app.init()

	if err := server.Serve(app.config.Connection.Port, app.routes()); err != nil {
		slog.Error("unable to serve application", err)
		os.Exit(1)
	}

	app.dbpool.Close()

	if err := app.tp.Shutdown(context.Background()); err != nil {
		slog.Error("error shutting down trace provider", err)
		os.Exit(1)
	}
}
