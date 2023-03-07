package middleware

import (
	"context"
	"crypto/sha256"
	"errors"
	"expvar"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/services"
	"github.com/seanflannery10/core/pkg/errs"
	"github.com/seanflannery10/core/pkg/helpers"
	"github.com/seanflannery10/core/pkg/validator"
	semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"golang.org/x/exp/slog"
)

func StartSpan(env services.Env) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			routePattern := chi.RouteContext(r.Context()).RoutePattern()
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

func Authenticate(env services.Env) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Vary", "Authorization")

			authorizationHeader := r.Header.Get("Authorization")

			if authorizationHeader == "" {
				ctx := context.WithValue(r.Context(), helpers.UserContextKey, data.AnonymousUser)
				next.ServeHTTP(w, r.WithContext(ctx))

				return
			}

			headerParts := strings.Split(authorizationHeader, " ")
			if len(headerParts) != 2 || headerParts[0] != "Bearer" {
				_ = render.Render(w, r, errs.ErrInvalidAuthenticationToken())
				return
			}

			token := headerParts[1]

			v := validator.New()

			v.Check(token != "", "token", "must be provided")
			v.Check(len(token) == 26, "token", "must be 26 bytes long")

			if v.HasErrors() {
				_ = render.Render(w, r, errs.ErrInvalidAuthenticationToken())
				return
			}

			tokenHash := sha256.Sum256([]byte(token))

			user, err := env.Queries.GetUserFromToken(r.Context(), data.GetUserFromTokenParams{
				Hash:   tokenHash[:],
				Scope:  data.ScopeAuthentication,
				Expiry: pgtype.Timestamptz{Time: time.Now(), Valid: true},
			})
			if err != nil {
				switch {
				case errors.Is(err, pgx.ErrNoRows):
					_ = render.Render(w, r, errs.ErrInvalidAuthenticationToken())
				default:
					_ = render.Render(w, r, errs.ErrServerError(err))
				}

				return
			}

			r = r.WithContext(context.WithValue(r.Context(), helpers.UserContextKey, user))

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := helpers.ContextGetUser(r)

		if user.IsAnonymous() {
			_ = render.Render(w, r, errs.ErrAuthenticationRequired)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Metrics(next http.Handler) http.Handler {
	totalRequestsReceived := expvar.NewInt("total_requests_received")
	totalResponsesSent := expvar.NewInt("total_responses_sent")
	totalProcessingTimeMicroseconds := expvar.NewInt("total_processing_time_μs")
	totalResponsesSentByStatus := expvar.NewMap("total_responses_sent_by_status")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		metrics := httpsnoop.CaptureMetrics(next, w, r)

		totalRequestsReceived.Add(1)
		totalResponsesSent.Add(1)
		totalProcessingTimeMicroseconds.Add(metrics.Duration.Microseconds())
		totalResponsesSentByStatus.Add(strconv.Itoa(metrics.Code), 1)
	})
}

func RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler { //nolint:goerr113
					panic(rvr)
				}

				slog.Log(r.Context(), slog.LevelError, "panic recovery error", "error", rvr, "stack", string(debug.Stack()))

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}

				render.JSON(w, r, &errs.ErrResponse{
					Message: "the server encountered a problem and could not process your json",
				})
			}
		}()

		next.ServeHTTP(w, r)
	})
}
