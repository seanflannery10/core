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

//	@title			Core API
//	@version		0.1.0
//	@description	This is a core server.

//	@contact.name	Core Support
//	@contact.url	github.com/seanflannery10/
//	@contact.email	seanflannery10@gmail.com

// @host		api.seanflannery.dev
// @BasePath	/v1
func main() {
	app := &application{}

	app.init()

	if err := server.Serve(app.env.Config.Port, app.routes()); err != nil {
		slog.Error("unable to serve application", err)
		os.Exit(exitError)
	}

	app.shutdown()
}
