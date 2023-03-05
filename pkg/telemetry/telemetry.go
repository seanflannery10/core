package telemetry

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/seanflannery10/core/pkg/helpers"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func New(endpoint, env string) (*sdktrace.TracerProvider, error) {
	exp, err := otlptrace.New(context.Background(), otlptracegrpc.NewClient(
		otlptracegrpc.WithEndpoint(endpoint),
	))
	if err != nil {
		return nil, err
	}

	res, err := resource.Merge(resource.Default(), resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceVersionKey.String(helpers.GetVersion()),
		semconv.ServiceNameKey.String("core"),
		attribute.String("environment", env)))
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(1)),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return tp, nil
}

func TraceHandler(ctx context.Context, r *http.Request, tp *sdktrace.TracerProvider) oteltrace.Span {
	tracer := tp.Tracer("core")

	spanName := ""

	routePattern := chi.RouteContext(r.Context()).RoutePattern() //nolint:contextcheck

	if routePattern == "" {
		spanName = "/"
	} else {
		spanName = routePattern
	}

	spanName = r.Method + " " + spanName

	_, span := tracer.Start(ctx, spanName,
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
	)

	span.SetAttributes(semconv.HTTPRouteKey.String(routePattern))
	span.SetName(spanName)

	//// set status code attribute
	// span.SetAttributes(semconv.HTTPStatusCodeKey.Int(rrw.status))
	//
	//// set span status
	// spanStatus, spanMessage := semconv.SpanStatusFromHTTPStatusCode(rrw.status)
	// span.SetStatus(spanStatus, spanMessage)

	return span
}
