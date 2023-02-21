package main

import (
	"crypto/sha256"
	"errors"
	"expvar"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/felixge/httpsnoop"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/seanflannery10/core/internal/data"
	"github.com/seanflannery10/core/internal/helpers"
	"github.com/seanflannery10/core/internal/httperrors"
	"github.com/seanflannery10/core/internal/validator"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = helpers.ContextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			httperrors.InvalidAuthenticationToken(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); v.HasErrors() {
			httperrors.InvalidAuthenticationToken(w, r)
			return
		}

		tokenHash := sha256.Sum256([]byte(token))

		user, err := app.queries.GetUserFromToken(r.Context(), data.GetUserFromTokenParams{
			Hash:   tokenHash[:],
			Scope:  data.ScopePasswordReset,
			Expiry: pgtype.Timestamptz{Time: time.Now(), Valid: true},
		})
		if err != nil {
			switch {
			case errors.Is(err, pgx.ErrNoRows):
				httperrors.InvalidAuthenticationToken(w, r)
			default:
				httperrors.ServerError(w, r, err)
			}
			return
		}

		r = helpers.ContextSetUser(r, &user)

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := helpers.ContextGetUser(r)

		if user.IsAnonymous() {
			httperrors.AuthenticationRequired(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) metrics(next http.Handler) http.Handler {
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

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				err, _ := rvr.(error)
				if !errors.Is(err, http.ErrAbortHandler) {
					httperrors.ServerError(w, r, err)
				}
			}
		}()

		next.ServeHTTP(w, r)
	})
}
