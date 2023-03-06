package main

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/internal/services"
	"github.com/seanflannery10/core/pkg/server"
	"github.com/seanflannery10/core/pkg/telemetry"
	"golang.org/x/exp/slog"
)

type application struct {
	config          Config
	tracerProviders telemetry.TracerProviders
	dbpool          *pgxpool.Pool
	env             *services.Env
}

func main() {
	app := &application{}

	app.init()

	if err := server.Serve(app.config.Port, app.routes()); err != nil {
		slog.Error("unable to serve application", err)
		os.Exit(1)
	}

	app.shutdown()
}
