package telemetry

import (
	"context"

	"github.com/seanflannery10/core/pkg/helpers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

type TracerProviders struct {
	Standard *sdktrace.TracerProvider
	Error    *sdktrace.TracerProvider
}

func NewTracerProviders(endpoint, env string) (TracerProviders, error) {
	exp, err := otlptrace.New(context.Background(), otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(endpoint),
	))
	if err != nil {
		return TracerProviders{}, err
	}

	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceVersionKey.String(helpers.GetVersion()),
		semconv.ServiceNameKey.String("core"),
		attribute.String("environment", env)))
	if err != nil {
		return TracerProviders{}, err
	}

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	telemetry := TracerProviders{}

	telemetry.Standard = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(1)),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	telemetry.Error = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	return telemetry, err
}
