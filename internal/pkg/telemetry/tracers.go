package telemetry

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

type Tracers struct {
	Standard oteltrace.Tracer
	Error    oteltrace.Tracer
}

func NewTrace(r *http.Request, tracer oteltrace.Tracer) oteltrace.Span {
	spanName := ""

	routePattern := chi.RouteContext(r.Context()).RoutePattern()

	if routePattern == "" {
		spanName = "/"
	} else {
		spanName = routePattern
	}

	spanName = r.Method + " " + spanName

	_, span := tracer.Start(r.Context(), spanName,
		oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
		oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
		oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest("core", routePattern, r)...),
		oteltrace.WithSpanKind(oteltrace.SpanKindServer),
	)

	span.SetAttributes(semconv.HTTPRouteKey.String(routePattern))
	span.SetName(spanName)

	return span
}
