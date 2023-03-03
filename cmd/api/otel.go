package main

import (
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/exp/slog"
)

func (app *application) otel() {
	client := otlptracegrpc.NewClient()

	exp, err := otlptrace.New(ctx, client)
	if err != nil {
		slog.Error("unable to initialize exporter", err)
		os.Exit(1)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
	)

	defer func() {
		_ = exp.Shutdown(ctx)
		_ = tp.Shutdown(ctx)
	}()

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
}
