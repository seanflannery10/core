package main

import (
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/seanflannery10/core/internal/pkg/server"
	"github.com/seanflannery10/core/internal/services"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/exp/slog"
)

type application struct {
	dbpool *pgxpool.Pool
	tp     *sdktrace.TracerProvider
	env    *services.Env
}

func main() {
	app := &application{}

	app.init()

	if err := server.Serve(app.env.Config.Port, app.routes()); err != nil {
		slog.Error("unable to serve application", err)
		os.Exit(exitError)
	}

	app.shutdown()
}
