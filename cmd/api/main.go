package main

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/pkg/mailer"
	"github.com/seanflannery10/core/pkg/server"
	"github.com/seanflannery10/core/pkg/telemetry"
	"golang.org/x/exp/slog"
)

type application struct {
	config          Config
	mailer          mailer.Mailer
	tracerProviders telemetry.TracerProviders
	dbpool          *pgxpool.Pool
}

func main() {
	app := &application{}

	app.init()

	if err := server.Serve(app.config.Connection.Port, app.routes()); err != nil {
		slog.Error("unable to serve application", err)
		os.Exit(1)
	}

	app.dbpool.Close()

	if err := app.tracerProviders.Standard.Shutdown(context.Background()); err != nil {
		slog.Error("error shutting down standard trace provider", err)
		os.Exit(1)
	}

	if err := app.tracerProviders.Error.Shutdown(context.Background()); err != nil {
		slog.Error("error shutting down error trace provider", err)
		os.Exit(1)
	}
}
