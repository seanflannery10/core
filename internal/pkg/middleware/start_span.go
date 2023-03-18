package middleware

import (
	"net/http"

	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/pkg/helpers"
	"github.com/seanflannery10/core/internal/services"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func StartSpan(env *services.Env) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routePattern := helpers.GetRoutePattern(r)
			spanName := routePattern

			if spanName == "" {
				spanName = "/"
			}

			spanName = r.Method + " " + spanName

			_, span := env.Tracer.Start(r.Context(), spanName,
				oteltrace.WithAttributes(semconv.NetAttributesFromHTTPRequest("tcp", r)...),
				oteltrace.WithAttributes(semconv.EndUserAttributesFromHTTPRequest(r)...),
				oteltrace.WithAttributes(semconv.HTTPServerAttributesFromHTTPRequest("core", routePattern, r)...),
				oteltrace.WithSpanKind(oteltrace.SpanKindServer),
			)

			span.SetAttributes(semconv.HTTPRouteKey.String(routePattern))
			span.SetName(spanName)

			r = r.WithContext(oteltrace.ContextWithSpan(r.Context(), span))

			next.ServeHTTP(w, r)
		})
	}
}

func LogUser(r *http.Request, user *data.User) *http.Request {
	span := oteltrace.SpanFromContext(r.Context())
	span.SetAttributes(attribute.Int64("user.id", user.ID))
	span.SetAttributes(attribute.String("user.name", user.Name))
	span.SetAttributes(attribute.String("user.email", user.Email))

	return r.WithContext(oteltrace.ContextWithSpan(r.Context(), span))
}
