package main

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

func (app *application) shutdown() {
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
