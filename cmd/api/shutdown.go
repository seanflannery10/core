package main

import (
	"context"
	"os"

	"golang.org/x/exp/slog"
)

func (app *application) shutdown() {
	app.dbpool.Close()

	if err := app.tp.Shutdown(context.Background()); err != nil {
		slog.Error("error shutting down trace provider", err)
		os.Exit(1) //nolint:revive
	}
}
